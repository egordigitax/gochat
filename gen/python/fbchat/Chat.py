# automatically generated by the FlatBuffers compiler, do not modify

# namespace: fbchat

import flatbuffers
from flatbuffers.compat import import_numpy
np = import_numpy()

class Chat(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAs(cls, buf, offset=0):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = Chat()
        x.Init(buf, n + offset)
        return x

    @classmethod
    def GetRootAsChat(cls, buf, offset=0):
        """This method is deprecated. Please switch to GetRootAs."""
        return cls.GetRootAs(buf, offset)
    # Chat
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

    # Chat
    def Title(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # Chat
    def UnreadCount(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int32Flags, o + self._tab.Pos)
        return 0

    # Chat
    def LastMessage(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(8))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # Chat
    def LastAuthor(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(10))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

    # Chat
    def MediaUrl(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        if o != 0:
            return self._tab.String(o + self._tab.Pos)
        return None

def ChatStart(builder):
    builder.StartObject(5)

def Start(builder):
    ChatStart(builder)

def ChatAddTitle(builder, title):
    builder.PrependUOffsetTRelativeSlot(0, flatbuffers.number_types.UOffsetTFlags.py_type(title), 0)

def AddTitle(builder, title):
    ChatAddTitle(builder, title)

def ChatAddUnreadCount(builder, unreadCount):
    builder.PrependInt32Slot(1, unreadCount, 0)

def AddUnreadCount(builder, unreadCount):
    ChatAddUnreadCount(builder, unreadCount)

def ChatAddLastMessage(builder, lastMessage):
    builder.PrependUOffsetTRelativeSlot(2, flatbuffers.number_types.UOffsetTFlags.py_type(lastMessage), 0)

def AddLastMessage(builder, lastMessage):
    ChatAddLastMessage(builder, lastMessage)

def ChatAddLastAuthor(builder, lastAuthor):
    builder.PrependUOffsetTRelativeSlot(3, flatbuffers.number_types.UOffsetTFlags.py_type(lastAuthor), 0)

def AddLastAuthor(builder, lastAuthor):
    ChatAddLastAuthor(builder, lastAuthor)

def ChatAddMediaUrl(builder, mediaUrl):
    builder.PrependUOffsetTRelativeSlot(4, flatbuffers.number_types.UOffsetTFlags.py_type(mediaUrl), 0)

def AddMediaUrl(builder, mediaUrl):
    ChatAddMediaUrl(builder, mediaUrl)

def ChatEnd(builder):
    return builder.EndObject()

def End(builder):
    return ChatEnd(builder)


class ChatT(object):

    # ChatT
    def __init__(self):
        self.title = None  # type: str
        self.unreadCount = 0  # type: int
        self.lastMessage = None  # type: str
        self.lastAuthor = None  # type: str
        self.mediaUrl = None  # type: str

    @classmethod
    def InitFromBuf(cls, buf, pos):
        chat = Chat()
        chat.Init(buf, pos)
        return cls.InitFromObj(chat)

    @classmethod
    def InitFromPackedBuf(cls, buf, pos=0):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, pos)
        return cls.InitFromBuf(buf, pos+n)

    @classmethod
    def InitFromObj(cls, chat):
        x = ChatT()
        x._UnPack(chat)
        return x

    # ChatT
    def _UnPack(self, chat):
        if chat is None:
            return
        self.title = chat.Title()
        self.unreadCount = chat.UnreadCount()
        self.lastMessage = chat.LastMessage()
        self.lastAuthor = chat.LastAuthor()
        self.mediaUrl = chat.MediaUrl()

    # ChatT
    def Pack(self, builder):
        if self.title is not None:
            title = builder.CreateString(self.title)
        if self.lastMessage is not None:
            lastMessage = builder.CreateString(self.lastMessage)
        if self.lastAuthor is not None:
            lastAuthor = builder.CreateString(self.lastAuthor)
        if self.mediaUrl is not None:
            mediaUrl = builder.CreateString(self.mediaUrl)
        ChatStart(builder)
        if self.title is not None:
            ChatAddTitle(builder, title)
        ChatAddUnreadCount(builder, self.unreadCount)
        if self.lastMessage is not None:
            ChatAddLastMessage(builder, lastMessage)
        if self.lastAuthor is not None:
            ChatAddLastAuthor(builder, lastAuthor)
        if self.mediaUrl is not None:
            ChatAddMediaUrl(builder, mediaUrl)
        chat = ChatEnd(builder)
        return chat
