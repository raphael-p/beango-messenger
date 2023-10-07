// Event handler, emits "send-message" if enter is pressed with no modifier key
const sendMessageOnEnter = (event) => {
    if (event.key === "Enter" && !(event.shiftKey || event.altKey || event.ctrlKey || event.metaKey)) {
        event.preventDefault();
        htmx.trigger(event.target, "send-message");
    }
};

// Event handler, scroll to the bottom if at top or less than half the viewport away from the bottom
const adjustScroll = (event) => {
    const element = event.currentTarget;
    const distanceFromBottom = element.scrollTopMax - element.scrollTop;
    const viewportHeight = element.scrollHeight - element.scrollTopMax;
    if (element.scrollTop === 0 || distanceFromBottom < viewportHeight * 0.5 ) {
        element.scrollTop = element.scrollTopMax; 
    }
} 
