// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../BTNGGoldToken.sol";
import "../BTNGCustody.sol";
import "../BTNGGoldOracle.sol";

/**
 * @title BTNG Gold System Tests
 * @dev Comprehensive test suite for the BTNG Gold smart contract system
 */
contract BTNGGoldSystemTest is Test {
    BTNGGoldToken public goldToken;
    BTNGCustody public custody;
    BTNGGoldOracle public oracle;

    address public owner = address(1);
    address public custodian = address(2);
    address public user = address(3);
    address public feeder = address(4);

    uint256 public constant INITIAL_GOLD_PRICE = 6000000000; // $60 per gram
    uint256 public constant MINT_AMOUNT = 100 * 10**18; // 100 BTNG
    uint256 public constant GOLD_GRAMS = 100; // 100 grams

    function setUp() public {
        vm.startPrank(owner);

        // Deploy Oracle
        oracle = new BTNGGoldOracle();

        // Deploy Custody with placeholder token address
        custody = new BTNGCustody(address(0), address(oracle));

        // Deploy Token
        goldToken = new BTNGGoldToken(address(custody));

        // Update Custody with real token address
        custody = new BTNGCustody(address(goldToken), address(oracle));

        // Update Token with new custody address
        goldToken.updateCustodyContract(address(custody));

        // Setup roles
        custody.addCustodian(custodian);
        oracle.addFeeder(feeder);

        vm.stopPrank();
    }

    function testDeployment() public {
        assertEq(goldToken.name(), "BTNG Gold Token");
        assertEq(goldToken.symbol(), "BTNG");
        assertEq(address(goldToken.custodyContract()), address(custody));
        assertEq(address(custody.oracle()), address(oracle));
        assertTrue(custody.isCustodian(custodian));
        assertTrue(oracle.isFeeder(feeder));
    }

    function testGoldPriceUpdates() public {
        vm.prank(feeder);
        oracle.updateGoldPrice(6500000000); // $65 per gram

        assertEq(oracle.getGoldPrice(), 6500000000);
    }

    function testCountryReserveUpdates() public {
        vm.prank(feeder);
        oracle.updateCountryReserve("GH", 10000000); // Ghana: 10 tonnes

        (uint256 reserves, uint256 lastUpdated, bool isActive) = oracle.getCountryReserve("GH");
        assertEq(reserves, 10000000);
        assertTrue(isActive);
        assertGt(lastUpdated, 0);
    }

    function testSovereignValueCalculation() public {
        vm.prank(feeder);
        oracle.updateCountryReserve("GH", 8700000); // Ghana: 8.7 tonnes

        uint256 sovereignValue = oracle.calculateSovereignValue("GH", 31072940); // Ghana population
        uint256 expectedValue = (8700000 * 1000000 * 1000000) / 31072940; // (grams * 1e6) / population

        assertEq(sovereignValue, expectedValue);
    }

    function testGoldDepositAndMint() public {
        // Deposit gold
        vm.prank(custodian);
        custody.depositGold("BATCH001", GOLD_GRAMS, "ipfs://certificate-hash");

        // Request mint
        vm.prank(user);
        custody.requestMint(MINT_AMOUNT, "BATCH001");

        // Execute mint (via oracle)
        vm.prank(address(oracle));
        custody.executeMint(user, MINT_AMOUNT, "BATCH001");

        // Verify minting
        assertEq(goldToken.balanceOf(user), MINT_AMOUNT);
        assertEq(goldToken.getGoldBacking(user), GOLD_GRAMS);
        assertTrue(goldToken.verifyGoldBacking());
    }

    function testGoldRedemption() public {
        // First mint some tokens
        vm.prank(custodian);
        custody.depositGold("BATCH002", GOLD_GRAMS, "ipfs://certificate-hash");

        vm.prank(user);
        custody.requestMint(MINT_AMOUNT, "BATCH002");

        vm.prank(address(oracle));
        custody.executeMint(user, MINT_AMOUNT, "BATCH002");

        // Now test redemption
        vm.prank(user);
        custody.requestRedemption(MINT_AMOUNT, "REDEEM001", "123 Gold Street, Accra");

        vm.prank(custodian);
        custody.executeRedemption("REDEEM001");

        // Verify redemption
        assertEq(goldToken.balanceOf(user), 0);
        assertEq(goldToken.getGoldBacking(user), 0);
        assertTrue(goldToken.verifyGoldBacking());
    }

    function testTokenTransfers() public {
        // Mint tokens to user
        vm.prank(custodian);
        custody.depositGold("BATCH003", GOLD_GRAMS, "ipfs://certificate-hash");

        vm.prank(user);
        custody.requestMint(MINT_AMOUNT, "BATCH003");

        vm.prank(address(oracle));
        custody.executeMint(user, MINT_AMOUNT, "BATCH003");

        address recipient = address(5);
        uint256 transferAmount = 50 * 10**18; // 50 BTNG

        // Transfer tokens
        vm.prank(user);
        goldToken.transfer(recipient, transferAmount);

        // Verify transfer
        assertEq(goldToken.balanceOf(user), MINT_AMOUNT - transferAmount);
        assertEq(goldToken.balanceOf(recipient), transferAmount);
        assertEq(goldToken.getGoldBacking(user), GOLD_GRAMS - 50);
        assertEq(goldToken.getGoldBacking(recipient), 50);
    }

    function testEmergencyPause() public {
        vm.prank(owner);
        goldToken.pause();

        // Try to transfer (should fail)
        vm.prank(user);
        vm.expectRevert("Pausable: paused");
        goldToken.transfer(address(5), 100);
    }

    function testAccessControl() public {
        // Non-owner trying to add custodian
        vm.prank(user);
        vm.expectRevert("Ownable: caller is not the owner");
        custody.addCustodian(address(6));

        // Non-feeder trying to update price
        vm.prank(user);
        vm.expectRevert("BTNGGoldOracle: caller is not authorized feeder");
        oracle.updateGoldPrice(6500000000);
    }

    function testVerificationSystem() public {
        bytes32 requestId = keccak256(abi.encodePacked("test-request"));

        vm.prank(feeder);
        oracle.requestVerification(requestId, 0, MINT_AMOUNT, "BATCH001");

        bool result = oracle.getVerificationResult(requestId);
        assertTrue(result);
    }

    function testTotalReserveCalculation() public {
        uint256 totalReserves = oracle.getTotalAfricanReserves();
        assertGt(totalReserves, 0); // Should have initial data
    }

    function testMaxSupplyLimit() public {
        uint256 maxSupply = goldToken.MAX_SUPPLY();

        // Try to mint more than max supply
        vm.prank(custodian);
        custody.depositGold("BATCH004", maxSupply / 10**18 + 1, "ipfs://certificate-hash");

        vm.prank(user);
        custody.requestMint(maxSupply + 10**18, "BATCH004");

        vm.prank(address(oracle));
        vm.expectRevert("BTNGGoldToken: exceeds max supply");
        custody.executeMint(user, maxSupply + 10**18, "BATCH004");
    }

    function testInsufficientReserves() public {
        // Try to mint without sufficient gold reserves
        vm.prank(user);
        custody.requestMint(MINT_AMOUNT, "NONEXISTENT");

        vm.prank(address(oracle));
        vm.expectRevert("BTNGCustody: invalid gold batch");
        custody.executeMint(user, MINT_AMOUNT, "NONEXISTENT");
    }
}