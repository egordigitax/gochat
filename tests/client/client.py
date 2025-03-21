import sys
import os
import flatbuffers
import websocket

sys.path.insert(0, os.path.abspath("./gen/python/"))

import fbchat.GetMessageFromClientRequest as GetMessageFromClientRequest
import fbchat.RootMessage as RootMessage
import fbchat.ActionType as ActionType


def build_get_message_request(text: str) -> bytes:
    inner_builder = flatbuffers.Builder(256)
    text_offset = inner_builder.CreateString(text)

    GetMessageFromClientRequest.Start(inner_builder)
    GetMessageFromClientRequest.AddText(inner_builder, text_offset)
    msg_offset = GetMessageFromClientRequest.End(inner_builder)
    inner_builder.Finish(msg_offset)

    inner_bytes = inner_builder.Output()

    outer_builder = flatbuffers.Builder(512)
    payload_vector = outer_builder.CreateByteVector(inner_bytes)

    RootMessage.Start(outer_builder)
    RootMessage.AddActionType(outer_builder, ActionType.ActionType.GET_MESSAGE)
    RootMessage.AddPayload(outer_builder, payload_vector)
    root_offset = RootMessage.End(outer_builder)
    outer_builder.Finish(root_offset)

    return bytes(outer_builder.Output())


def send_ws_message():
    ws = websocket.create_connection(
        "ws://127.0.0.1:8081/fb/messages?chat_id=5f2e0cd8-87e3-4cc9-98ed-1fcd77a954de",
        header=[
            "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTE5MjlmOTMtZmQxNy00ZTlkLWIzOGMtMzFmNGMyNmZhNTFjIn0.7yLstLzY3KP7ZO-ZQU5_0hhgx_pf2Z8-DTZA8DhSe_8"
        ]
    )

    payload = build_get_message_request("privet")
    ws.send_binary(payload)
    response = ws.recv()
    print("Raw response:", response)

    ws.close()


if __name__ == "__main__":
    send_ws_message()
