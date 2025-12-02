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
  console.log(`Deploying LocalToken from ${wallet.address}`);

  const sources = {
    "ERC20.sol": {
      content: fs.readFileSync(path.join(__dirname, "../contracts/ERC20.sol"), "utf8"),
    },
    "LocalToken.sol": {
      content: fs.readFileSync(path.join(__dirname, "../contracts/LocalToken.sol"), "utf8"),
    },
  };

  const input = {
    language: "Solidity",
    sources,
    settings: {
      optimizer: {
        enabled: false,
        runs: 200,
      },
      outputSelection: {
        "*": {
          "*": ["abi", "evm.bytecode"],
        },
      },
    },
  };

  const output = JSON.parse(solc.compile(JSON.stringify(input)));
  if (output.errors && output.errors.length) {
    let hasErrors = false;
    for (const error of output.errors) {
      const message = error.formattedMessage || error.message;
      if (error.severity === "error") {
        hasErrors = true;
        console.error(message);
      } else {
        console.warn(message);
      }
    }
    if (hasErrors) {
      throw new Error("solc compilation failed");
    }
  }

  const contract = output.contracts["LocalToken.sol"].LocalToken;
  if (!contract) {
    throw new Error("LocalToken artifact missing from compiler output");
  }

  const { abi, evm } = contract;
  const bytecode = evm?.bytecode?.object;
  if (!bytecode) {
    throw new Error("Bytecode missing from compilation output");
  }

  const initialSupply = BigInt(process.env.INITIAL_SUPPLY || "0");
  console.log(`Deploying LocalToken with initial supply ${initialSupply.toString()} wei-style units`);

  const factory = new ethers.ContractFactory(abi, bytecode, wallet);
  const token = await factory.deploy(initialSupply);
  await token.waitForDeployment();
  const tokenAddress = await token.getAddress();
  console.log(`LocalToken deployed at ${tokenAddress}`);

  console.log("Done.");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
