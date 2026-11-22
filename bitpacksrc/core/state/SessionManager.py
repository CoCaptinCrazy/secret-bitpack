import uuid

from bitpack.core.state.StateManager import StateManager


class SessionManager:
    def __init__(self):
        self.state = StateManager()

    def start(self):
        session_id = str(uuid.uuid4())

        self.state.set(
            "current_session",
            session_id
        )

        return session_id

    def active(self):
        return self.state.get(
            "current_session"
        )

    def end(self):
        self.state.remove(
            "current_session"
        )