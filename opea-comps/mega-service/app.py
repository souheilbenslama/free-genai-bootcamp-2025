from attr import has
from comps.cores.mega.constants import ServiceType, ServiceRoleType
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo
)

from comps.cores.proto.docarray import LLMParams, RerankerParms, RetrieverParms
from comps import MicroService, ServiceOrchestrator
from fastapi.responses import StreamingResponse
import os 

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 7000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 9000)


class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):  # Change port to 6000
        print("hello")
        self.host = host
        self.port = port
        self.endpoint = "/v1/examples"  # Add endpoint attribute
        self.megaservice = ServiceOrchestrator()

    def add_remote_service(self):
        # Create the source service (our main service)
        embedding_service = MicroService(
            name="embedding",
            host=EMBEDDING_SERVICE_HOST_IP,
            port=EMBEDDING_SERVICE_PORT,
            service_type=ServiceType.EMBEDDING,
            use_remote_service=True,  # Mark as remote service
            endpoint="/v1/embeddings"  # Add endpoint for embeddings
        )
        
        # Create the LLM service
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        
        # Add both services to the orchestrator
        self.megaservice.add(embedding_service)
        self.megaservice.add(llm)
        
        # Set up the flow from source service to LLM
        self.megaservice.flow_to(embedding_service, llm)
    
    def start(self):

        self.service = MicroService(
            self.__class__.__name__,
            service_role=ServiceRoleType.MEGASERVICE,
            host=self.host,
            port=self.port,
            endpoint="/v1/examples",
            input_datatype=ChatCompletionRequest,
            output_datatype=ChatCompletionResponse,
        )

        self.service.add_route(self.endpoint, self.handle_request, methods=["POST"])
        print(f"Service configured with endpoint: {self.endpoint}")
        self.service.start()
        
    async def handle_request(self, request: ChatCompletionRequest) -> dict:
        """Handle incoming requests by forwarding them through the service flow"""
        try:
            # Create LLM parameters with defaults
            llm_params = LLMParams(
                max_tokens=getattr(request, 'max_tokens', 1024),
                top_k=getattr(request, 'top_k', 10),
                top_p=getattr(request, 'top_p', 0.95),
                temperature=getattr(request, 'temperature', 0.01),
                frequency_penalty=getattr(request, 'frequency_penalty', 0.0),
                presence_penalty=getattr(request, 'presence_penalty', 0.0),
                repetition_penalty=getattr(request, 'repetition_penalty', 1.03),
                stream=getattr(request, 'stream', False),
                chat_template=getattr(request, 'chat_template', None)
            )

            # Get messages
            messages = request.messages
            if not messages:
                return ChatCompletionResponse(
                    model="example-model",
                    choices=[ChatCompletionResponseChoice(
                        index=0,
                        message=ChatMessage(role="assistant", content="No messages provided"),
                        finish_reason="stop"
                    )],
                    usage=UsageInfo()
                ).dict()

            # Schedule the processing through the service flow
            result_dict, runtime_graph = await self.megaservice.schedule(
                initial_inputs={
                    "messages": messages,
                    "text": messages[-1].content
                },
                llm_parameters=llm_params
            )

            # Handle streaming response if present
            for node, response in result_dict.items():
                if isinstance(response, StreamingResponse):
                    return response

            # Get the final response from the last node
            last_node = runtime_graph.all_leaves()[-1]
            response = result_dict[last_node].get("text", "No response generated")

            # Create and return the response
            return ChatCompletionResponse(
                model="example-model",
                choices=[ChatCompletionResponseChoice(
                    index=0,
                    message=ChatMessage(role="assistant", content=response),
                    finish_reason="stop"
                )],
                usage=UsageInfo()
            ).dict()

        except Exception as e:
            print(f"Error processing request: {str(e)}")
            return ChatCompletionResponse(
                model="example-model",
                choices=[ChatCompletionResponseChoice(
                    index=0,
                    message=ChatMessage(
                        role="assistant",
                        content=f"Error processing request: {str(e)}"
                    ),
                    finish_reason="stop"  # Changed from 'error' to 'stop'
                )],
                usage=UsageInfo()
            ).dict()

example = ExampleService()
example.add_remote_service()
example.start()
