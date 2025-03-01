import chromadb
import os

# setup Chroma in-memory, for easy prototyping. Can add persistence easily!
client = chromadb.Client()

# Create collection. get_collection, get_or_create_collection, delete_collection also available!
collection = client.create_collection("all-my-documents")

def read_documents_from_directory(directory: str):
    documents = []
    metadatas = []
    ids = []
    for filename in os.listdir(directory):
        if filename.endswith('.txt'):
            with open(os.path.join(directory, filename), 'r') as file:
                documents.append(file.read())
                metadatas.append({"source": filename})
                ids.append(filename)
    return documents, metadatas, ids

# Read documents from the directory
documents, metadatas, ids = read_documents_from_directory(os.path.dirname(__file__))

# Add docs to the collection
collection.add(
    documents=documents,
    metadatas=metadatas,
    ids=ids
)

# Query/search 2 most similar results. You can also .get by id
results = collection.query(
    query_texts=["This is a query document"],
    n_results=2,
    # where={"metadata_field": "is_equal_to_this"}, # optional filter
    # where_document={"$contains":"search_string"}  # optional filter
)