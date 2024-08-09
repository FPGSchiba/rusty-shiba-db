supportedOS=("linux" "darwin" "windows")
supportedArches=("arm64" "amd64")

for os in ${supportedOS[@]}; do
  for arch in ${supportedArches[@]}; do
    extionsion=""
    if [ $os == 'windows' ]; then extionsion='.exe'; fi
    env GOOS="$os" GOARCH="$arch" go build -o "./rusty-shiba-db-$os-$arch$extionsion" rsdb/src
  done
done