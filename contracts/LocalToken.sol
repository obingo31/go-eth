// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./ERC20.sol";

/// @title LocalToken - Minimal concrete contract for deploying the Solady-style ERC20 base
contract LocalToken is ERC20 {
    constructor(uint256 initialSupply) {
        _mint(msg.sender, initialSupply);
    }

    function name() public view override returns (string memory) {
        return "Local Token";
    }

    function symbol() public view override returns (string memory) {
        return "LOCAL";
    }
}
