// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title DemoToken - a minimal ERC-20 token for local testing
contract DemoToken {
    string public constant name = "Demo Token";
    string public constant symbol = "DEMO";
    uint8 public constant decimals = 18;

    uint256 public totalSupply;
    address public owner;
    mapping(address => uint256) private balances;
    mapping(address => mapping(address => uint256)) private allowances;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(uint256 initialSupply) {
        totalSupply = initialSupply;
        owner = msg.sender;
        balances[msg.sender] = initialSupply;
        emit Transfer(address(0), msg.sender, initialSupply);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "DemoToken: caller is not the owner");
        _;
    }

    function balanceOf(address account) external view returns (uint256) {
        return balances[account];
    }

    function allowance(address tokenOwner, address spender) external view returns (uint256) {
        return allowances[tokenOwner][spender];
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        _transfer(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowances[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        uint256 currentAllowance = allowances[from][msg.sender];
        require(currentAllowance >= amount, "ERC20: allowance exceeded");
        allowances[from][msg.sender] = currentAllowance - amount;
        _transfer(from, to, amount);
        return true;
    }

    function mint(address to, uint256 amount) external onlyOwner returns (bool) {
        require(to != address(0), "DemoToken: mint to zero address");
        totalSupply += amount;
        balances[to] += amount;
        emit Transfer(address(0), to, amount);
        return true;
    }

    function _transfer(address from, address to, uint256 amount) internal {
        require(to != address(0), "ERC20: transfer to the zero address");
        uint256 fromBalance = balances[from];
        require(fromBalance >= amount, "ERC20: balance too low");
        balances[from] = fromBalance - amount;
        balances[to] += amount;
        emit Transfer(from, to, amount);
    }
}
