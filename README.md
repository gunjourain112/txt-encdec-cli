# txt-encdec-cli
My lightweight text encryption/decryption script that I use on Fedora.

리눅스 전용

### Build and Run
```bash
git clone https://github.com/gunjourain112/txt-encdec-cli.git
cd txt-encdec-cli
go clean -cache -modcache
go build -o enc .
./enc
```

### Alias Setting (Optional)
```bash
echo "alias enc='$(pwd)/enc'" >> ~/.bashrc
source ~/.bashrc
enc
```

### .Desktop File (Optional)
```bash
vi ~/.local/share/applications/encryptor.desktop
```

```bash
[Desktop Entry]
Version=1.0
Type=Application
Name=Encryptor
Comment=AES-GCM text enc dec tool
Exec=/home/?/appsh/txt-encdec-cli/enc
Icon=/home/?/appsh/txt-encdec-cli/enc.png
Terminal=true
Categories=Utility;Security;
```

```bash
chmod +x ~/.local/share/applications/encryptor.desktop
update-desktop-database ~/.local/share/applications
```
