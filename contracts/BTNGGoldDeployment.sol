// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./BTNGGoldToken.sol";
import "./BTNGCustody.sol";
import "./BTNGGoldOracle.sol";

/**
 * @title BTNG Gold System Deployment
 * @dev Deployment script for the complete BTNG Gold system
 */
contract BTNGGoldDeployment {
    BTNGGoldToken public goldToken;
    BTNGCustody public custody;
    BTNGGoldOracle public oracle;

    address public owner;
    address public oracleAddress;
    address public custodyAddress;
    address public tokenAddress;

    event SystemDeployed(
        address indexed token,
        address indexed custody,
        address indexed oracle,
        address owner
    );

    constructor() {
        owner = msg.sender;

        // Deploy Oracle first
        oracle = new BTNGGoldOracle();
        oracleAddress = address(oracle);

        // Deploy Custody with Oracle address
        custody = new BTNGCustody(address(0), oracleAddress); // Token address will be set later
        custodyAddress = address(custody);

        // Deploy Token with Custody address
        goldToken = new BTNGGoldToken(custodyAddress);
        tokenAddress = address(goldToken);

        // Update Custody contract with Token address
        custody = new BTNGCustody(tokenAddress, oracleAddress);
        custodyAddress = address(custody);

        // Update Token contract with new Custody address
        goldToken.updateCustodyContract(custodyAddress);

        // Transfer ownerships
        oracle.transferOwnership(owner);
        custody.transferOwnership(owner);
        goldToken.transferOwnership(owner);

        emit SystemDeployed(tokenAddress, custodyAddress, oracleAddress, owner);
    }

    /**
     * @dev Get all deployed contract addresses
     */
    function getContractAddresses()
        external
        view
        returns (
            address token,
            address custody_,
            address oracle_
        )
    {
        return (tokenAddress, custodyAddress, oracleAddress);
    }
}