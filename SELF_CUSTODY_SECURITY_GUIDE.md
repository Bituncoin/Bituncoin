# BTNG Self-Custody Security Guide

## Overview

BTNG implements a revolutionary self-custody system where users maintain complete control of their assets through offline-first, hardware-isolated private keys. This guide explains how to securely manage your BTNG wealth without relying on third parties.

## The Ghost-Key System

### Key Architecture

BTNG uses Multi-Party Computation (MPC) to shard your private key:

- **Primary Shard**: Stored encrypted on your device (never leaves your control)
- **Recovery Shard**: Held by BTNG network (cannot reconstruct key alone)
- **Backup Shards**: Optional additional shards for enhanced security

### Shard Generation

When you create your 35-digit BTNG address:

1. **MPC Ceremony**: Your device participates in a secure key generation ceremony
2. **Primary Shard**: Immediately encrypted and stored locally
3. **Recovery Setup**: Recovery shard is created and distributed across the network
4. **Verification**: All shards are cryptographically verified before activation

## Managing Your Primary Shard

### Device Storage

```typescript
// Primary shard is stored in hardware-backed secure storage
const primaryShard = await deviceSecureStorage.get('btng-primary-shard');

// Encrypted with device biometric authentication
const decryptedShard = await biometricDecrypt(primaryShard, userBiometric);
```

### Backup Strategies

#### Option 1: Hardware Backup

- Store encrypted shard on dedicated hardware security module (HSM)
- Never connect HSM to internet
- Use for recovery only

#### Option 2: Manual Backup

- Generate recovery phrases from shard
- Store in multiple secure locations
- Never store digitally in cloud services

#### Option 3: Social Recovery

- Share recovery shards with trusted individuals
- Use Shamir's Secret Sharing for threshold recovery
- Minimum 3-of-5 shard requirement

## Transaction Signing Process

### Offline Preparation

1. **Online Phase**: Private Banker prepares transaction details
2. **Air-Gapped Transfer**: Transaction data moved to offline device
3. **Hardware Signing**: Private key shard used in isolated environment

### Signing Flow

```typescript
// Transaction prepared online
const preparedTx = await banker.prepareTransaction({
  amount: 50.00,
  recipient: 'BTNGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
  memo: 'Gold token purchase'
});

// Air-gapped transfer (QR code, NFC, etc.)
const offlineTx = transferToOfflineDevice(preparedTx);

// Hardware signing
const signedTx = await hardwareSign(offlineTx, primaryShard);

// Broadcast signed transaction
await network.broadcastTransaction(signedTx);
```

## Security Best Practices

### Device Security

- **Biometric Authentication**: Required for all shard access
- **Device Encryption**: Full disk encryption enabled
- **Remote Wipe**: Ability to remotely disable shards if device lost
- **Regular Updates**: Keep device firmware and BTNG app updated

### Network Security

- **No Internet Access**: Primary shard never touches internet-connected devices
- **Isolated Signing**: Use dedicated offline device for signing
- **Cold Storage**: Consider complete offline storage for large holdings

### Recovery Planning

- **Test Recovery**: Regularly test your recovery process
- **Update Contacts**: Keep recovery shard holders updated
- **Inheritance Planning**: Document recovery procedures for beneficiaries

## Risk Mitigation

### Loss of Device

1. **Immediate Action**: Disable lost device via recovery process
2. **Recovery Initiation**: Use backup shards to reconstruct key
3. **New Device Setup**: Generate new primary shard for replacement device

### Compromised Recovery

1. **Detection**: Monitor for unauthorized recovery attempts
2. **Invalidation**: Immediately invalidate compromised recovery shards
3. **Reconstruction**: Generate new recovery shards with trusted parties

### Network Attacks

1. **Offline Protection**: Since keys are offline, network attacks cannot access funds
2. **MPC Security**: No single party can reconstruct your private key
3. **Watchtower Monitoring**: Network continuously monitors for unusual activity

## Advanced Features

### Multi-Device Support

- Synchronize primary shard across multiple trusted devices
- Require quorum of devices for high-value transactions
- Automatic shard rotation for enhanced security

### Institutional Custody

- For organizations: Distributed shard storage across board members
- Legal compliance: Audit trails for all shard operations
- Insurance integration: Coverage for MPC-related incidents

### Emergency Access

- **Dead Man's Switch**: Automatic recovery after inactivity period
- **Trusted Third Party**: Optional integration with legal entities
- **Court-Ordered Recovery**: Legal mechanisms for asset recovery

## Troubleshooting

### Cannot Access Primary Shard

- Verify biometric authentication is working
- Check device storage integrity
- Use recovery process if primary access lost

### Transaction Signing Fails

- Ensure offline device is properly isolated
- Verify shard integrity with checksums
- Check hardware signing module functionality

### Recovery Process Issues

- Confirm all recovery shards are accessible
- Verify shard checksums match
- Contact BTNG support for MPC recovery assistance

## Support and Resources

### Documentation

- [MPC Technical Specification](./docs/mpc-spec.md)
- [Hardware Security Guide](./docs/hardware-security.md)
- [Recovery Procedures](./docs/recovery-guide.md)

### Community Support

- BTNG Security Forum
- 24/7 Emergency Recovery Line
- Local Sovereign Node Support

### Professional Services

- Security Audit Services
- Custody Consulting
- Institutional Setup Assistance

---

**Remember**: With great sovereignty comes great responsibility. Your assets are now truly yours - manage them accordingly.

*This guide is for informational purposes. Always consult with security professionals for your specific situation.*
