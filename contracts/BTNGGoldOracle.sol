// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title BTNG Gold Oracle
 * @dev Oracle contract for gold reserve data and verification
 * Provides real-time gold pricing and reserve verification
 */
contract BTNGGoldOracle is Ownable, Pausable, ReentrancyGuard {
    // Events
    event GoldPriceUpdated(uint256 oldPrice, uint256 newPrice, uint256 timestamp);
    event ReserveDataUpdated(string countryCode, uint256 newReserves);
    event OracleFeederAdded(address indexed feeder);
    event OracleFeederRemoved(address indexed feeder);
    event VerificationCompleted(bytes32 requestId, bool isValid);

    // Structs
    struct GoldPrice {
        uint256 price; // Price per gram in USD (with 8 decimals)
        uint256 timestamp;
        address updater;
    }

    struct CountryReserve {
        string countryCode;
        uint256 goldTonnes; // Gold reserves in tonnes
        uint256 lastUpdated;
        address updater;
        bool isActive;
    }

    // State variables
    mapping(address => bool) public authorizedFeeders;
    mapping(string => CountryReserve) public countryReserves;
    GoldPrice public latestGoldPrice;

    uint256 public constant PRICE_DECIMALS = 8;
    uint256 public constant MAX_PRICE_AGE = 24 hours; // Maximum age for price data
    uint256 public constant MAX_RESERVE_AGE = 7 days; // Maximum age for reserve data

    string[] public activeCountries;

    // Request tracking for verification
    mapping(bytes32 => bool) public verificationRequests;
    mapping(bytes32 => bool) public verificationResults;

    // Modifiers
    modifier onlyFeeder() {
        require(authorizedFeeders[msg.sender] || msg.sender == owner(), "BTNGGoldOracle: caller is not authorized feeder");
        _;
    }

    modifier validPrice(uint256 price) {
        require(price > 0, "BTNGGoldOracle: invalid price");
        require(price < 100000 * 10**PRICE_DECIMALS, "BTNGGoldOracle: price too high"); // Max $100k per gram
        _;
    }

    modifier validCountryCode(string calldata countryCode) {
        require(bytes(countryCode).length == 2, "BTNGGoldOracle: invalid country code");
        _;
    }

    /**
     * @dev Constructor
     */
    constructor() Ownable(msg.sender) {
        // Add deployer as initial feeder
        authorizedFeeders[msg.sender] = true;

        // Initialize with current known data (as of Feb 2026)
        _initializeCountryData();
    }

    /**
     * @dev Initialize country reserve data
     */
    function _initializeCountryData() internal {
        // Major African gold reserve holders
        _updateCountryReserve("DZ", 173600000); // Algeria: 173.6 tonnes
        _updateCountryReserve("EG", 79300000);  // Egypt: 79.3 tonnes
        _updateCountryReserve("LY", 116600000); // Libya: 116.6 tonnes
        _updateCountryReserve("ZA", 125300000); // South Africa: 125.3 tonnes
        _updateCountryReserve("GH", 8700000);   // Ghana: 8.7 tonnes
        _updateCountryReserve("NG", 21400000);  // Nigeria: 21.4 tonnes
        _updateCountryReserve("MA", 22100000);  // Morocco: 22.1 tonnes
        _updateCountryReserve("TN", 6800000);   // Tunisia: 6.8 tonnes

        // Set initial gold price (approx $60 per gram as of 2026)
        latestGoldPrice = GoldPrice({
            price: 6000000000, // $60.00 per gram
            timestamp: block.timestamp,
            updater: msg.sender
        });
    }

    /**
     * @dev Update gold price per gram
     * @param newPrice Price in USD with 8 decimals (e.g., 6000000000 = $60.00)
     */
    function updateGoldPrice(uint256 newPrice)
        external
        onlyFeeder
        whenNotPaused
        validPrice(newPrice)
        nonReentrant
    {
        uint256 oldPrice = latestGoldPrice.price;
        latestGoldPrice = GoldPrice({
            price: newPrice,
            timestamp: block.timestamp,
            updater: msg.sender
        });

        emit GoldPriceUpdated(oldPrice, newPrice, block.timestamp);
    }

    /**
     * @dev Update country gold reserves
     * @param countryCode ISO 3166-1 alpha-2 country code
     * @param goldTonnes Gold reserves in tonnes (with 6 decimals for precision)
     */
    function updateCountryReserve(
        string calldata countryCode,
        uint256 goldTonnes
    )
        external
        onlyFeeder
        whenNotPaused
        validCountryCode(countryCode)
        nonReentrant
    {
        _updateCountryReserve(countryCode, goldTonnes);
        emit ReserveDataUpdated(countryCode, goldTonnes);
    }

    /**
     * @dev Internal function to update country reserve
     */
    function _updateCountryReserve(string memory countryCode, uint256 goldTonnes) internal {
        if (!countryReserves[countryCode].isActive) {
            activeCountries.push(countryCode);
            countryReserves[countryCode].isActive = true;
        }

        countryReserves[countryCode] = CountryReserve({
            countryCode: countryCode,
            goldTonnes: goldTonnes,
            lastUpdated: block.timestamp,
            updater: msg.sender,
            isActive: true
        });
    }

    /**
     * @dev Get total African gold reserves
     * @return Total gold reserves in tonnes
     */
    function getTotalAfricanReserves() external view returns (uint256) {
        uint256 total = 0;
        for (uint256 i = 0; i < activeCountries.length; i++) {
            total += countryReserves[activeCountries[i]].goldTonnes;
        }
        return total;
    }

    /**
     * @dev Get gold price in USD per gram
     * @return Price with 8 decimals
     */
    function getGoldPrice() external view returns (uint256) {
        require(block.timestamp - latestGoldPrice.timestamp <= MAX_PRICE_AGE, "BTNGGoldOracle: gold price too old");
        return latestGoldPrice.price;
    }

    /**
     * @dev Get country reserve data
     * @param countryCode ISO country code
     * @return goldTonnes The gold reserves in tonnes
     * @return lastUpdated The last update timestamp
     * @return isActive Whether the country is active
     */
    function getCountryReserve(string calldata countryCode)
        external
        view
        returns (uint256 goldTonnes, uint256 lastUpdated, bool isActive)
    {
        CountryReserve memory reserve = countryReserves[countryCode];
        require(reserve.isActive, "BTNGGoldOracle: country not found");

        return (reserve.goldTonnes, reserve.lastUpdated, reserve.isActive);
    }

    /**
     * @dev Calculate sovereign value for a country
     * @param countryCode ISO country code
     * @param population Country population
     * @return Sovereign value per person in BTNG
     */
    function calculateSovereignValue(
        string calldata countryCode,
        uint256 population
    )
        external
        view
        returns (uint256)
    {
        CountryReserve memory reserve = countryReserves[countryCode];
        require(reserve.isActive, "BTNGGoldOracle: country not found");
        require(population > 0, "BTNGGoldOracle: invalid population");

        // Sovereign value = (Gold reserves in grams * 1e6) / Population
        uint256 goldGrams = reserve.goldTonnes * 1000000; // Convert tonnes to grams
        return (goldGrams * 1000000) / population; // Multiply by 1e6 for precision
    }

    /**
     * @dev Request verification of gold operation
     * @param requestId Unique request identifier
     * @param operationType Type of operation (0=mint, 1=redeem)
     * @param amount Amount involved
     * @param goldBatchId Gold batch identifier
     */
    function requestVerification(
        bytes32 requestId,
        uint8 operationType,
        uint256 amount,
        string calldata goldBatchId
    )
        external
        onlyFeeder
        whenNotPaused
    {
        require(!verificationRequests[requestId], "BTNGGoldOracle: request already exists");

        verificationRequests[requestId] = true;

        // Basic validation logic (can be extended)
        bool isValid = _validateOperation(operationType, amount, goldBatchId);

        verificationResults[requestId] = isValid;

        emit VerificationCompleted(requestId, isValid);
    }

    /**
     * @dev Internal validation logic
     */
    function _validateOperation(
        uint8 operationType,
        uint256 amount,
        string calldata goldBatchId
    )
        internal
        view
        returns (bool)
    {
        // Basic validations
        if (amount == 0) return false;

        // Check if price data is fresh
        if (block.timestamp - latestGoldPrice.timestamp > MAX_PRICE_AGE) {
            return false;
        }

        // For minting operations, ensure we have sufficient reserves
        if (operationType == 0) { // Mint
            uint256 totalReserves = 0;
            for (uint256 i = 0; i < activeCountries.length; i++) {
                totalReserves += countryReserves[activeCountries[i]].goldTonnes;
            }
            uint256 requiredGrams = amount / 10**18;
            if (requiredGrams > totalReserves) return false;
        }

        return true;
    }

    /**
     * @dev Get verification result
     * @param requestId Request identifier
     * @return Verification result
     */
    function getVerificationResult(bytes32 requestId) external view returns (bool) {
        require(verificationRequests[requestId], "BTNGGoldOracle: request not found");
        return verificationResults[requestId];
    }

    /**
     * @dev Add authorized feeder
     * @param feeder Address to authorize
     */
    function addFeeder(address feeder) external onlyOwner {
        require(feeder != address(0), "BTNGGoldOracle: invalid feeder address");
        require(!authorizedFeeders[feeder], "BTNGGoldOracle: already feeder");

        authorizedFeeders[feeder] = true;
        emit OracleFeederAdded(feeder);
    }

    /**
     * @dev Remove authorized feeder
     * @param feeder Address to remove
     */
    function removeFeeder(address feeder) external onlyOwner {
        require(authorizedFeeders[feeder], "BTNGGoldOracle: feeder not authorized");

        authorizedFeeders[feeder] = false;
        emit OracleFeederRemoved(feeder);
    }

    /**
     * @dev Emergency pause
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @dev Unpause after emergency
     */
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Get active countries count
     */
    function getActiveCountriesCount() external view returns (uint256) {
        return activeCountries.length;
    }

    /**
     * @dev Get active country by index
     */
    function getActiveCountry(uint256 index) external view returns (string memory) {
        require(index < activeCountries.length, "BTNGGoldOracle: index out of bounds");
        return activeCountries[index];
    }

    /**
     * @dev Check if feeder is authorized
     */
    function isFeeder(address account) external view returns (bool) {
        return authorizedFeeders[account];
    }
}