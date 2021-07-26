# list-addons-in-win-programs

## How to build the native messaging host and its installer?

On Windows 10 + WSL:

1. [Install and setup Golang](https://golang.org/doc/install) on your Linux environment.
2. Install go-msi https://github.com/mh-cbon/go-msi *via an MSI to your Windows environment*.
3. Install WiX Toolset https://wixtoolset.org/releases/ to your Windows environment.
4. Set PATH to go-msi (ex. `C:\Program Files\go-msi`) and WiX Toolse (ex. `C:\Program Files (x86)\WiX Toolset v3.11\bin`).
5. Run `make host`.
   Then `.exe` files and a batch file to build MSI will be generated.
6. Double-click the generated `webextensions\native-messaging-host\build_msi.bat` on your Windows environment.
   Then two MSIs will be generated.

