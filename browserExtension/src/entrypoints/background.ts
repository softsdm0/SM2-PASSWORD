import { usePasswordStore } from "@/stores/password";

export default defineBackground(() => {
  console.log("Hello background!", { id: browser.runtime.id });

  browser.runtime.onMessage.addListener((message, _, sendResponse) => {
    console.log("Content script received message:", message);
    if (message.type === "getPassword") {
      // const password = await usePasswordStore().getPassword(message.id);
      // console.log("getPassword", password);
      // sendResponse(password);
      usePasswordStore().getPassword(message.id).then(sendResponse);
    }
    return true;
  });
});
