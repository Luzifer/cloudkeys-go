import GibberishAES from 'gibberish-aes/src/gibberish-aes.js'
import { SHA256 } from 'sha2'

const opensslDecrypt = (plaintext, password, callback) => {
  callback(GibberishAES.dec(plaintext, password), null)
}

const opensslEncrypt = (ciphertext, password, callback) => {
  callback(GibberishAES.enc(ciphertext, password), null)
}

const sha256sum = (plaintext, callback) => {
  callback(SHA256(plaintext).toString('hex'))
}

window.opensslDecrypt = opensslDecrypt
window.opensslEncrypt = opensslEncrypt
window.sha256sum = sha256sum
