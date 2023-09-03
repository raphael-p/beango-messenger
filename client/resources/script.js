const sendMessageEvent = new Event("send-message");
const handleKeypress = (event) => {
    if (event.key === "Enter") {
        event.preventDefault();
        event.target.dispatchEvent(sendMessageEvent);
        // event.target.value = "";
    }
};

const handleChatOpened = () => {
    for (const textarea of document.getElementsByClassName("message-input")) {
        textarea.addEventListener("keypress", handleKeypress);
    }
}
document.addEventListener("chat-opened", () => {
    handleChatOpened();
});