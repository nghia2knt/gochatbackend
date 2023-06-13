from pymongo import MongoClient
from pymongo.errors import PyMongoError
from robot.api.deco import keyword
from robot.libraries.BuiltIn import BuiltIn


class MongoDBLibrary:
    def __init__(self, connection_string):
        try:
            self.client = MongoClient(connection_string)
            self.db = None
        except PyMongoError as e:
            raise AssertionError(f"Failed to connect to MongoDB: {str(e)}")

    @keyword('Connect to MongoDB')
    def connect_to_mongodb(self, db_name):
        try:
            self.db = self.client[db_name]
            BuiltIn().log(f'Connected to MongoDB: {self.client.host}:{self.client.port} - Database: {db_name}')
            print("success to connect mongo db")
        except PyMongoError as e:
            raise AssertionError(f"Failed to connect to MongoDB: {str(e)}")

    @keyword('Execute MongoDB Query')
    def execute_mongodb_query(self, collection_name, query):
        try:
            collection = self.db[collection_name]
            result = list(collection.find(query))
            BuiltIn().log(f'Result: {result}')
            return result
        except PyMongoError as e:
            raise AssertionError(f"Failed to execute MongoDB query: {str(e)}")
        
    @keyword('Insert One To MongoDB')
    def insert_one_to_mongodb(self, collection_name, record):
        try:
            collection = self.db[collection_name]
            result = collection.insert_one(record)
            BuiltIn().log(f'Result: {result.inserted_id}')
            return result
        except PyMongoError as e:
            raise AssertionError(f"Failed to execute MongoDB command for insert: {str(e)}")
