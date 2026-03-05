export const encryptPassword = (password: string): string => {
  // TODO: 先直接返回，不进行加密，二期再做
  // return CryptoJS.AES.encrypt(password, SECRET_KEY).toString()
  return password
}

export const decryptPassword = (encryptedPassword: string): string => {
  // TODO: 先直接返回，不进行解密，二期再做
  // const bytes = CryptoJS.AES.decrypt(encryptedPassword, SECRET_KEY)
  // return bytes.toString(CryptoJS.enc.Utf8)
  return encryptedPassword
}
