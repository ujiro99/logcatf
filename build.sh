VERSION=$(git describe --tags)
gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}" -os="windows darwin linux" -ldflags="-s -w -X main.version=${VERSION}"
cd dist
mv  darwin_386_logcatf        logcatf
zip darwin_386_logcatf        logcatf     -qm
mv  darwin_amd64_logcatf      logcatf
zip darwin_amd64_logcatf      logcatf     -qm
mv  linux_386_logcatf         logcatf
zip linux_386_logcatf         logcatf     -qm
mv  linux_amd64_logcatf       logcatf
zip linux_amd64_logcatf       logcatf     -qm
mv  linux_arm_logcatf         logcatf
zip linux_arm_logcatf         logcatf     -qm
mv  windows_386_logcatf.exe   logcatf.exe
zip windows_386_logcatf       logcatf.exe -qm
mv  windows_amd64_logcatf.exe logcatf.exe
zip windows_amd64_logcatf     logcatf.exe -qm
cd -
