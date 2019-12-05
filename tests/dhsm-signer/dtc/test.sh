#!/usr/bin/env sh
set -e
# Sign zone
echo "Signing zone example.com..."
./dhsm-signer sign -p ./dtc.so -f ./example.com -3 -z example.com -o example.com.signed -c || echo "Cannot sign properly" && exit

echo "Signing successful!"

# Verify the zone already signed
echo "Verifying previous signature..."
./dhsm-signer verify -f ./example.com.signed || echo "Cannot verify signature :(" && exit

echo "Verification successful!"

# Reset HSM keys
echo "Resetting keys..."
./dhsm-signer reset-keys -p ./dtc.so || echo "Cannot reset keys :(" && exit

echo "Reset successful! you can kill this process with ^C."