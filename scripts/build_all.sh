echo "Please run me at root dir"

if [ -d "output" ]; then
    echo "Output dir existed, deleting and recreating..."
    rm -rf output
fi
mkdir -p output/services

pushd backend/go-services || exit

for i in *; do
  cd services
  if [ "$i" != "health" ]; then
    echo "$i"
    capName="${name}"
    cd "$i" || exit
    go build -o "../../../output/services/$i/${capName}Service"
    cd ..
  fi
done

popd || exit

mkdir -p output/gateway

cd backend/api-gateway || exit

go build -o "../../output/gateway/Gateway"

echo "OK!"



