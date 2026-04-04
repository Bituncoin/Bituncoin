using System;
using System.Collections.Generic;
using System.Text.Json;

namespace BTNGWallet
{
  public class BTNGWalletService
  {
    private readonly Dictionary<string, decimal> _ledger = new();
    private readonly List<object> _transactionLog = new();

    public string Mint(string wallet, decimal amount)
    {
      if (string.IsNullOrWhiteSpace(wallet))
        throw new ArgumentException("Wallet address is required", nameof(wallet));
      if (amount <= 0)
        throw new ArgumentException("Mint amount must be positive", nameof(amount));

      _ledger[wallet] = _ledger.GetValueOrDefault(wallet, 0m) + amount;

      var tx = new
      {
        type = "mint",
        wallet,
        amount,
        balance_after = _ledger[wallet],
        timestamp = DateTime.UtcNow.ToString("O")
      };
      _transactionLog.Add(tx);
      return JsonSerializer.Serialize(tx, new JsonSerializerOptions { WriteIndented = true });
    }

    public string Melt(string wallet, decimal amount)
    {
      if (string.IsNullOrWhiteSpace(wallet))
        throw new ArgumentException("Wallet address is required", nameof(wallet));
      if (amount <= 0)
        throw new ArgumentException("Melt amount must be positive", nameof(amount));

      var currentBalance = _ledger.GetValueOrDefault(wallet, 0m);
      if (currentBalance < amount)
        throw new InvalidOperationException(
          $"Insufficient BTNG: wallet {wallet} has {currentBalance}, requested {amount}");

      _ledger[wallet] -= amount;

      var tx = new
      {
        type = "melt",
        wallet,
        amount,
        balance_after = _ledger[wallet],
        timestamp = DateTime.UtcNow.ToString("O")
      };
      _transactionLog.Add(tx);
      return JsonSerializer.Serialize(tx, new JsonSerializerOptions { WriteIndented = true });
    }

    public decimal GetBalance(string wallet)
    {
      return _ledger.GetValueOrDefault(wallet, 0m);
    }

    public string GetLedger()
    {
      return JsonSerializer.Serialize(_ledger, new JsonSerializerOptions { WriteIndented = true });
    }

    public string GetTransactionLog()
    {
      return JsonSerializer.Serialize(_transactionLog, new JsonSerializerOptions { WriteIndented = true });
    }

    public string GetBalanceSummary()
    {
      var response = new
      {
        statement_id = Guid.NewGuid().ToString(),
        period = new
        {
          start = "2026-01-01T00:00:00Z",
          end = "2026-01-31T23:59:59Z"
        },
        global_summary = new
        {
          starting_balance = new { btng = "1200.50", usd = "84210.33" },
          ending_balance = new { btng = "1450.75", usd = "101233.90" },
          net_change = new { btng = "250.25", usd = "16923.57" },
          total_inflows = new { btng = "500.00", usd = "33800.00" },
          total_outflows = new { btng = "249.75", usd = "16876.43" }
        },
        chain_breakdown = new
        {
          btng = new { starting = "800.00", ending = "950.00", net_change = "150.00" },
          bitcoin = new { starting = "0.25", ending = "0.30", net_change = "0.05" },
          ethereum = new { starting = "3.00", ending = "4.20", net_change = "1.20" }
        },
        address_breakdown = new List<object>
                {
                    new
                    {
                        address = "0x12ab...89ff",
                        label = "Main ETH Wallet",
                        chain = "ethereum",
                        starting_balance = new { native = "1.50", btng = "300.00", usd = "20250.00" },
                        ending_balance = new { native = "2.10", btng = "420.00", usd = "28350.00" },
                        net_change = new { native = "0.60", btng = "120.00", usd = "8100.00" }
                    }
                }
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string AddAddress(string chain, string type, string label, string addressOrXpub)
    {
      var response = new
      {
        address_id = Guid.NewGuid().ToString(),
        chain = chain,
        type = type,
        label = label,
        address_or_xpub = addressOrXpub,
        imported_addresses = new List<string> { "bc1qxyz...", "bc1qabc..." },
        status = "added"
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string ListAddresses()
    {
      var response = new
      {
        addresses = new List<object>
                {
                    new
                    {
                        address_id = Guid.NewGuid().ToString(),
                        chain = "btng",
                        type = "address",
                        label = "BTNG Main",
                        address = "btng1qxyz...",
                        transaction_count = 42
                    },
                    new
                    {
                        address_id = Guid.NewGuid().ToString(),
                        chain = "bitcoin",
                        type = "xpub",
                        label = "BTC XPUB",
                        address_or_xpub = "xpub6CUGRU...",
                        derived_addresses = 12,
                        transaction_count = 128
                    }
                }
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string ImportFromBlockchair()
    {
      var response = new
      {
        source = "blockchair_wow",
        imported = new List<object>
                {
                    new { chain = "bitcoin", address = "bc1qxyz...", label = "BTC Wallet 1" },
                    new { chain = "ethereum", address = "0xabc123...", label = "ETH Wallet 1" }
                },
        status = "success"
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string SetPriceAlert(string currency, string condition, decimal threshold, string notificationType = "email")
    {
      var response = new
      {
        alert_id = Guid.NewGuid().ToString(),
        currency = currency,
        condition = condition,
        threshold = threshold,
        notification_type = notificationType,
        created_at = DateTime.UtcNow.ToString("O"),
        status = "active",
        message = $"Alert set for {currency} price {condition} {threshold}"
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string GetPriceAlerts()
    {
      var response = new
      {
        alerts = new List<object>
                {
                    new
                    {
                        alert_id = Guid.NewGuid().ToString(),
                        currency = "USD",
                        condition = "above",
                        threshold = 150.00m,
                        notification_type = "email",
                        created_at = "2026-03-29T10:00:00Z",
                        status = "active",
                        last_triggered = (string)null
                    },
                    new
                    {
                        alert_id = Guid.NewGuid().ToString(),
                        currency = "GHS",
                        condition = "below",
                        threshold = 1600.00m,
                        notification_type = "push",
                        created_at = "2026-03-29T11:30:00Z",
                        status = "active",
                        last_triggered = "2026-03-29T14:15:00Z"
                    }
                },
        total_active = 2
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }

    public string RemovePriceAlert(string alertId)
    {
      var response = new
      {
        alert_id = alertId,
        status = "removed",
        removed_at = DateTime.UtcNow.ToString("O"),
        message = "Price alert successfully removed"
      };
      return JsonSerializer.Serialize(response, new JsonSerializerOptions { WriteIndented = true });
    }
  }
}
