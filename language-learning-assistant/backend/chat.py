import boto3
import streamlit as st
from typing import Optional, Dict, Any
from dotenv import load_dotenv
from mistralai import Mistral

import os
# Model ID
load_dotenv()
api_key = os.environ["MISTRAL_API_KEY"]
model = "mistral-large-latest"

class MistralChat:
    def __init__(self, model_id: str = model, api_key: str = api_key):
        """Initialize Bedrock chat client"""
        self.mistral_client = Mistral(api_key=api_key)
        self.model_id = model_id

    def generate_response(self, message: str, inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using Amazon Bedrock"""
        if inference_config is None:
            inference_config = {"temperature": 0.7}

        messages = [{
            "role": "user",
            "content":  message
        }]

        try:
            response = self.mistral_client.chat.complete(
                model=self.model_id,
                messages=messages,
                temperature=inference_config["temperature"]
            )
            return response
            
        except Exception as e:
            print(f"Error generating response: {str(e)}")  # Changed from st.error to print for CLI
            return None

def run_chat():
    st.title("Chat Application")
    chat = MistralChat()

    # Initialize chat history in session state
    if "messages" not in st.session_state:
        st.session_state.messages = []

    # Display chat history
    for message in st.session_state.messages:
        with st.chat_message(message["role"]):
            st.markdown(message["content"])

    # Chat input
    if prompt := st.chat_input("What's your message?"):
        # Add user message to chat history
        st.session_state.messages.append({"role": "user", "content": prompt})
        with st.chat_message("user"):
            st.markdown(prompt)

        # Get bot response
        with st.chat_message("assistant"):
            response = chat.generate_response(prompt)
            if response:
                st.markdown(response.choices[0].message.content)
                st.session_state.messages.append({"role": "assistant", "content": response.choices[0].message.content})

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1:
        # Run as Streamlit app if arguments are provided
        run_chat()
    else:
        # Run as CLI if no arguments
        chat = MistralChat()
        print("Chat initialized. Type '/exit' to end conversation.")
        while True:
            user_input = input("\nYou: ")
            if user_input.lower() == '/exit':
                break
            response = chat.generate_response(user_input)
            if response:
                print("Bot:", response.choices[0].message.content)
