# Mepm
A simple cross-platform password manager written in **Go** using the **Fyne** GUI framework.
It works on all major platforms and has been tested on **Linux** and **Android**.

The app uses **SQLite** for storage, **Argon2** for key derivation, and **AES** encryption to keep your passwords secure.

> ⚠️ **Warning**
> This app is under development and may never be fully finished.
> It was never meant to be a real-world password manager — it’s a testing ground and hobby project to learn more about **Go** and **encryption**.
## Screenshots





<p>
  <img src="screenshots/android.png" height="300">
  <img src="screenshots/desktop.png" height="300">
</p>


App is under heavy development, and it may never be finished !!!
## Todos

Security:
- [x] encryption / decryption with AES
- [x] key generation with argon2
- [x] salt generation

Db:
- [x] created models for password and note
- [x] insert, update, delete for password
- [ ] insert, update, delete for note

Gui:
- [x] password insert screen
- [x] password edit screen
- [x] delete password
- [ ] note insert screen
- [ ] note edit screen
- [ ] delete note
- [ ] additional setting screen
