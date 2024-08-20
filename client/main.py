# use example_pb2_grpc to call the servic

import grpc
import example_pb2
import example_pb2_grpc
import time

def run():
    channel = grpc.insecure_channel('localhost:8080')
    stub = example_pb2_grpc.GreeterStub(channel)
    response = stub.SayHello(example_pb2.HelloRequest(name='you'))
    print("Greeter client received: " + response.message)

if __name__ == '__main__':
    run()