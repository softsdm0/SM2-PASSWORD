import useStoredValue from "@/utils/useStoredValue";
import { LocalPasswordServerUrl } from "@/stores/storages";
import { usePasswordStore } from "@/stores/password";
import { IPasswordInfo } from "@/types/password";
export default defineContentScript({
  matches: ["*://*/*"],
  async main(ctx) {
    // 检查是否需要显示已有的账号和密码
    document.addEventListener("click", (e) => {
      if (e.target instanceof HTMLInputElement) {
        checkWhetherShowUserAndPassword();
      }
    });
  },
});

async function checkWhetherShowUserAndPassword() {
  console.log(window.location.hostname);
  // 如果不在列表中存在
  const mathList = await checkHostNameWhetherMatch(window.location.hostname);
  if (mathList.length === 0) {
    console.log("hostname, 不在列表中存在");
    return;
  }

  // 获取页面中所有密码输入框
  const passwordInputs = document.querySelectorAll("input");

  let usernameInput: HTMLInputElement | null = null;
  let passwordInput: HTMLInputElement | null = null;

  // 获取到type=password的输入框和他上一个账号的输入框。
  for (let i = 0; i < passwordInputs.length; i++) {
    const input = passwordInputs[i];
    if (input.type === "password") {
      passwordInput = input;
      usernameInput = passwordInputs[i - 1] as HTMLInputElement;
    }
  }
  if (!usernameInput && !passwordInput) {
    return;
  }

  // 在usernameInput标签添加一个下拉框，用于显示匹配到的密码名称
  const div = document.createElement("div");
  div.style.position = "absolute";
  div.style.zIndex = "9999";
  div.style.backgroundColor = "white";
  div.style.border = "1px solid #ddd";
  div.style.borderRadius = "4px";
  div.style.boxShadow = "0 2px 5px rgba(0,0,0,0.2)";
  div.style.maxHeight = "200px";
  div.style.overflowY = "auto";

  // 添加点击外部关闭功能
  const handleClickOutside = (event: MouseEvent) => {
    if (!div.contains(event.target as Node)) {
      div.remove();
      document.removeEventListener("click", handleClickOutside);
    }
  };

  // 添加到usernameInput的后面
  usernameInput?.parentNode?.insertBefore(div, usernameInput.nextSibling);

  // 添加事件监听
  setTimeout(() => {
    document.addEventListener("click", handleClickOutside);
  }, 0);

  // 定位下拉框
  if (usernameInput) {
    const rect = usernameInput.getBoundingClientRect();

    // 紧贴输入框右侧
    div.style.top = `${rect.height}px`;
    div.style.left = `${rect.width + 5}px`; // 右侧5px间距
  }
  // 给select添加内容
  for (const password of mathList) {
    const option = document.createElement("div"); // 改为div元素更合适
    option.style.padding = "8px";
    option.style.cursor = "pointer";
    option.textContent = `${password.AppName} ${password.Account}`;

    // 添加点击事件, 点击后填充账号和密码
    option.addEventListener("click", async () => {
      if (usernameInput) {
        usernameInput.focus(); // 自动聚焦
        // 填充账号
        usernameInput.value = password.Account;

        // 获取密码
        if (passwordInput) {
          // 通过浏览器插件传递消息给background.ts进行处理, 不然会跨域
          const response = await browser.runtime.sendMessage({
            type: "getPassword",
            id: password.ID,
          });
          if (response) {
            passwordInput.focus(); // 自动聚焦
            // 填充密码
            passwordInput!.value = response.Password;
          }
        }

        // 填充后移除下拉框
        div.remove();
        document.removeEventListener("click", handleClickOutside);
      }
    });

    // 添加悬停效果
    option.addEventListener("mouseenter", () => {
      option.style.backgroundColor = "#f5f5f5";
    });
    option.addEventListener("mouseleave", () => {
      option.style.backgroundColor = "transparent";
    });

    div.appendChild(option);
  }
}

// 看看域名是否在密码服务器中存在
async function checkHostNameWhetherMatch(
  hostname: string
): Promise<IPasswordInfo[]> {
  const mathList: IPasswordInfo[] = [];
  const passwordList = await usePasswordStore().getPasswordList();
  if (passwordList) {
    for (const password of passwordList) {
      // console.log(`hostname: ${hostname}, server: ${password.Url}`);
      if (password.Url.includes(hostname)) {
        console.log("匹配到了");
        mathList.push(password);
      }
    }
  }
  return mathList;
}
