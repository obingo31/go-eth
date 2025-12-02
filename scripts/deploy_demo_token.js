#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const solc = require("solc");
const { ethers } = require("ethers");

async function main() {
  const rpcUrl = process.env.RPC_URL || "http://127.0.0.1:8545";
  const privateKey = process.env.PRIVATE_KEY;
  if (!privateKey) {
    throw new Error("Set PRIVATE_KEY (hex string) in the environment");
  }

  const provider = new ethers.JsonRpcProvider(rpcUrl);
  const wallet = new ethers.Wallet(privateKey, provider);
  console.log(`Deploying from ${wallet.address}`);

  const source = fs.readFileSync(path.join(__dirname, "../contracts/DemoToken.sol"), "utf8");
  const input = {
    language: "Solidity",
    sources: {
      "DemoToken.sol": {
        content: source,
      },
    },
    settings: {
      optimizer: {
        enabled: false,
        runs: 200,
      },
      evmVersion: "london",
      outputSelection: {
        "*": {
          "*": ["abi", "evm.bytecode"],
        },
      },
    },
  };

  const output = JSON.parse(solc.compile(JSON.stringify(input)));
  if (output.errors && output.errors.length) {
    for (const error of output.errors) {
      console.error(error.formattedMessage || error.message);
    }
    throw new Error("solc compilation failed");
  }

  const contract = output.contracts["DemoToken.sol"].DemoToken;
  const abi = contract.abi;
  const bytecode = contract.evm.bytecode.object;
  if (!bytecode) {
    throw new Error("Bytecode missing from compilation output");
  }

  const initialSupply = BigInt(process.env.INITIAL_SUPPLY || "0");
  console.log(`Deploying DemoToken with initial supply ${initialSupply.toString()} wei-style units`);

  const factory = new ethers.ContractFactory(abi, bytecode, wallet);
  const demoToken = await factory.deploy(initialSupply);
  await demoToken.waitForDeployment();
  const tokenAddress = await demoToken.getAddress();
  console.log(`DemoToken deployed at ${tokenAddress}`);

  const mintTo = process.env.MINT_TO;
  const mintAmount = process.env.MINT_AMOUNT;
  if (mintTo && mintAmount) {
    console.log(`Minting ${mintAmount} tokens to ${mintTo}`);
    const tx = await demoToken.mint(mintTo, BigInt(mintAmount));
    const receipt = await tx.wait();
    console.log(`Mint tx hash: ${receipt.hash}`);
  }

  console.log("Done.");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
