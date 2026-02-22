export interface CountryData {
  name: string
  flag: string
  capital: string
  currency: string
  currencyCode: string
  population: string
  mobileMoney: string[]
  status: 'active' | 'launching' | 'planned'
  trustNodes: number
  activeUsers: number
  description: string
  highlights: string[]
}

export const countryData: Record<string, CountryData> = {
  ghana: {
    name: 'Ghana',
    flag: '🇬🇭',
    capital: 'Accra',
    currency: 'Ghanaian Cedi',
    currencyCode: 'GHS',
    population: '33M',
    mobileMoney: ['MTN Mobile Money', 'Vodafone Cash', 'AirtelTigo Money'],
    status: 'active',
    trustNodes: 12,
    activeUsers: 145000,
    description: 'Ghana is a founding member of the BTNG Trust Union, leading West African financial inclusion through sovereign identity infrastructure.',
    highlights: [
      'First West African Trust Union node',
      'Integrated with 3 major mobile money providers',
      'Gold Card acceptance across 200+ merchants',
      'Cross-border remittance corridors established'
    ]
  },
  kenya: {
    name: 'Kenya',
    flag: '🇰🇪',
    capital: 'Nairobi',
    currency: 'Kenyan Shilling',
    currencyCode: 'KES',
    population: '55M',
    mobileMoney: ['M-Pesa', 'Airtel Money'],
    status: 'active',
    trustNodes: 18,
    activeUsers: 320000,
    description: 'Kenya pioneered mobile money in Africa and now leads the BTNG Trust Union expansion in East Africa with the largest user base.',
    highlights: [
      'Largest BTNG user base in Africa',
      'M-Pesa full integration operational',
      '18 Trust Union nodes across regions',
      'Government partnership for national ID integration'
    ]
  },
  nigeria: {
    name: 'Nigeria',
    flag: '🇳🇬',
    capital: 'Abuja',
    currency: 'Nigerian Naira',
    currencyCode: 'NGN',
    population: '220M',
    mobileMoney: ['OPay', 'PalmPay', 'Kuda', 'MTN MoMo'],
    status: 'launching',
    trustNodes: 24,
    activeUsers: 89000,
    description: 'As Africa\'s largest economy, Nigeria\'s BTNG integration connects 220M citizens to sovereign identity and digital value infrastructure.',
    highlights: [
      'Africa\'s largest economy integration',
      '24 Trust Union nodes deployment',
      'Multiple fintech partnerships active',
      'Lagos financial hub as regional anchor'
    ]
  },
  togo: {
    name: 'Togo',
    flag: '🇹🇬',
    capital: 'Lomé',
    currency: 'West African CFA franc',
    currencyCode: 'XOF',
    population: '8.8M',
    mobileMoney: ['Togocom Mobile Money', 'Moov Money'],
    status: 'active',
    trustNodes: 4,
    activeUsers: 42000,
    description: 'Togo brings CFA franc integration to the BTNG network, creating bridges between West African monetary unions and sovereign identity systems.',
    highlights: [
      'CFA franc integration complete',
      'Regional remittance hub',
      'Government digital ID partnership',
      'Cross-border corridor with Ghana'
    ]
  },
  uganda: {
    name: 'Uganda',
    flag: '🇺🇬',
    capital: 'Kampala',
    currency: 'Ugandan Shilling',
    currencyCode: 'UGX',
    population: '48M',
    mobileMoney: ['MTN Mobile Money', 'Airtel Money'],
    status: 'active',
    trustNodes: 10,
    activeUsers: 125000,
    description: 'Uganda\'s strong mobile money adoption makes it an ideal Trust Union partner, with high penetration in rural and urban areas.',
    highlights: [
      'High mobile money penetration',
      'Rural financial inclusion leader',
      'Agricultural sector integration',
      'Cross-border trade with Kenya active'
    ]
  },
  'ivory-coast': {
    name: 'Ivory Coast (Côte d\'Ivoire)',
    flag: '🇨🇮',
    capital: 'Yamoussoukro',
    currency: 'West African CFA franc',
    currencyCode: 'XOF',
    population: '28M',
    mobileMoney: ['Orange Money', 'MTN Mobile Money', 'Moov Money'],
    status: 'launching',
    trustNodes: 8,
    activeUsers: 34000,
    description: 'Ivory Coast is expanding BTNG Trust Union presence in Francophone West Africa with strong fintech ecosystem support.',
    highlights: [
      'Francophone West Africa anchor',
      'Orange Money integration in progress',
      'Abidjan tech hub participation',
      'CFA franc corridor expansion'
    ]
  },
  'burkina-faso': {
    name: 'Burkina Faso',
    flag: '🇧🇫',
    capital: 'Ouagadougou',
    currency: 'West African CFA franc',
    currencyCode: 'XOF',
    population: '22M',
    mobileMoney: ['Orange Money', 'Coris Money'],
    status: 'planned',
    trustNodes: 3,
    activeUsers: 8500,
    description: 'Burkina Faso\'s BTNG integration focuses on financial inclusion for underbanked populations and agricultural value chains.',
    highlights: [
      'Agricultural sector focus',
      'Financial inclusion priority',
      'Regional cooperation framework',
      'Pilot program in Ouagadougou'
    ]
  },
  'south-africa': {
    name: 'South Africa',
    flag: '🇿🇦',
    capital: 'Pretoria',
    currency: 'South African Rand',
    currencyCode: 'ZAR',
    population: '60M',
    mobileMoney: ['Vodacom M-Pesa', 'MTN MoMo', 'Standard Bank Instant Money'],
    status: 'launching',
    trustNodes: 15,
    activeUsers: 67000,
    description: 'South Africa anchors BTNG expansion in Southern Africa, leveraging advanced financial infrastructure and regional economic leadership.',
    highlights: [
      'Southern Africa regional hub',
      'Advanced banking system integration',
      'SADC cross-border framework',
      'Johannesburg as financial center'
    ]
  }
}
