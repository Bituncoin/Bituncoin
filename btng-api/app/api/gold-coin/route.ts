import { NextResponse } from 'next/server'

// African countries data with gold reserves and population
const africanCountriesData = [
  { name: 'Algeria', code: 'DZ', flag: '🇩🇿', goldReserves: 173.6, population: 43851044 },
  { name: 'Angola', code: 'AO', flag: '🇦🇴', goldReserves: 17.8, population: 32866272 },
  { name: 'Benin', code: 'BJ', flag: '🇧🇯', goldReserves: 0.2, population: 12123200 },
  { name: 'Botswana', code: 'BW', flag: '🇧🇼', goldReserves: 0.1, population: 2351627 },
  { name: 'Burkina Faso', code: 'BF', flag: '🇧🇫', goldReserves: 0.1, population: 20903273 },
  { name: 'Burundi', code: 'BI', flag: '🇧🇮', goldReserves: 0.1, population: 11890784 },
  { name: 'Cameroon', code: 'CM', flag: '🇨🇲', goldReserves: 0.1, population: 26545863 },
  { name: 'Cape Verde', code: 'CV', flag: '🇨🇻', goldReserves: 0.1, population: 555987 },
  { name: 'Central African Republic', code: 'CF', flag: '🇨🇫', goldReserves: 0.1, population: 4829767 },
  { name: 'Chad', code: 'TD', flag: '🇹🇩', goldReserves: 0.1, population: 16425864 },
  { name: 'Comoros', code: 'KM', flag: '🇰🇲', goldReserves: 0.1, population: 869601 },
  { name: 'Congo', code: 'CG', flag: '🇨🇬', goldReserves: 0.1, population: 5518087 },
  { name: 'Democratic Republic of the Congo', code: 'CD', flag: '🇨🇩', goldReserves: 0.1, population: 89561403 },
  { name: 'Djibouti', code: 'DJ', flag: '🇩🇯', goldReserves: 0.1, population: 988000 },
  { name: 'Egypt', code: 'EG', flag: '🇪🇬', goldReserves: 79.3, population: 102334404 },
  { name: 'Equatorial Guinea', code: 'GQ', flag: '🇬🇶', goldReserves: 0.1, population: 1402985 },
  { name: 'Eritrea', code: 'ER', flag: '🇪🇷', goldReserves: 0.1, population: 3546421 },
  { name: 'Eswatini', code: 'SZ', flag: '🇸🇿', goldReserves: 0.1, population: 1160164 },
  { name: 'Ethiopia', code: 'ET', flag: '🇪🇹', goldReserves: 0.1, population: 114963588 },
  { name: 'Gabon', code: 'GA', flag: '🇬🇦', goldReserves: 0.1, population: 2225728 },
  { name: 'Gambia', code: 'GM', flag: '🇬🇲', goldReserves: 0.1, population: 2416668 },
  { name: 'Ghana', code: 'GH', flag: '🇬🇭', goldReserves: 8.7, population: 31072940 },
  { name: 'Guinea', code: 'GN', flag: '🇬🇳', goldReserves: 0.1, population: 13132795 },
  { name: 'Guinea-Bissau', code: 'GW', flag: '🇬🇼', goldReserves: 0.1, population: 1968001 },
  { name: 'Ivory Coast', code: 'CI', flag: '🇨🇮', goldReserves: 0.1, population: 26378274 },
  { name: 'Kenya', code: 'KE', flag: '🇰🇪', goldReserves: 0.1, population: 53771296 },
  { name: 'Lesotho', code: 'LS', flag: '🇱🇸', goldReserves: 0.1, population: 2142249 },
  { name: 'Liberia', code: 'LR', flag: '🇱🇷', goldReserves: 0.1, population: 5057681 },
  { name: 'Libya', code: 'LY', flag: '🇱🇾', goldReserves: 116.6, population: 6871292 },
  { name: 'Madagascar', code: 'MG', flag: '🇲🇬', goldReserves: 0.1, population: 27691018 },
  { name: 'Malawi', code: 'MW', flag: '🇲🇼', goldReserves: 0.1, population: 19129952 },
  { name: 'Mali', code: 'ML', flag: '🇲🇱', goldReserves: 0.1, population: 20250833 },
  { name: 'Mauritania', code: 'MR', flag: '🇲🇷', goldReserves: 0.1, population: 4649658 },
  { name: 'Mauritius', code: 'MU', flag: '🇲🇺', goldReserves: 0.1, population: 1271768 },
  { name: 'Morocco', code: 'MA', flag: '🇲🇦', goldReserves: 22.1, population: 36910560 },
  { name: 'Mozambique', code: 'MZ', flag: '🇲🇿', goldReserves: 0.1, population: 31255435 },
  { name: 'Namibia', code: 'NA', flag: '🇳🇦', goldReserves: 0.1, population: 2540905 },
  { name: 'Niger', code: 'NE', flag: '🇳🇪', goldReserves: 0.1, population: 24206644 },
  { name: 'Nigeria', code: 'NG', flag: '🇳🇬', goldReserves: 21.4, population: 218541000 },
  { name: 'Rwanda', code: 'RW', flag: '🇷🇼', goldReserves: 0.1, population: 12952218 },
  { name: 'Sao Tome and Principe', code: 'ST', flag: '🇸🇹', goldReserves: 0.1, population: 219159 },
  { name: 'Senegal', code: 'SN', flag: '🇸🇳', goldReserves: 0.1, population: 16743927 },
  { name: 'Seychelles', code: 'SC', flag: '🇸🇨', goldReserves: 0.1, population: 98347 },
  { name: 'Sierra Leone', code: 'SL', flag: '🇸🇱', goldReserves: 0.1, population: 7976983 },
  { name: 'Somalia', code: 'SO', flag: '🇸🇴', goldReserves: 0.1, population: 15893222 },
  { name: 'South Africa', code: 'ZA', flag: '🇿🇦', goldReserves: 125.3, population: 59308690 },
  { name: 'South Sudan', code: 'SS', flag: '🇸🇸', goldReserves: 0.1, population: 11193725 },
  { name: 'Sudan', code: 'SD', flag: '🇸🇩', goldReserves: 0.1, population: 43849260 },
  { name: 'Tanzania', code: 'TZ', flag: '🇹🇿', goldReserves: 0.1, population: 59734218 },
  { name: 'Togo', code: 'TG', flag: '🇹🇬', goldReserves: 0.1, population: 8278724 },
  { name: 'Tunisia', code: 'TN', flag: '🇹🇳', goldReserves: 6.8, population: 11818619 },
  { name: 'Uganda', code: 'UG', flag: '🇺🇬', goldReserves: 0.1, population: 45741007 },
  { name: 'Zambia', code: 'ZM', flag: '🇿🇲', goldReserves: 0.1, population: 18383955 },
  { name: 'Zimbabwe', code: 'ZW', flag: '🇿🇼', goldReserves: 0.1, population: 14862924 }
];

// Calculate sovereign values
const countriesWithSovereignValue = africanCountriesData.map(country => ({
  ...country,
  sovereignValue: (country.goldReserves * 1000000) / country.population
}));

const totalGoldValue = countriesWithSovereignValue.reduce((sum, country) => sum + country.sovereignValue, 0);

export async function GET() {
  try {
    return NextResponse.json({
      success: true,
      data: {
        countries: countriesWithSovereignValue,
        totalGoldValue,
        totalCountries: countriesWithSovereignValue.length,
        totalPopulation: countriesWithSovereignValue.reduce((sum, c) => sum + c.population, 0),
        totalGoldReserves: countriesWithSovereignValue.reduce((sum, c) => sum + c.goldReserves, 0),
        lastUpdated: new Date().toISOString()
      }
    });
  } catch (error) {
    return NextResponse.json(
      { success: false, error: 'Failed to fetch gold coin data' },
      { status: 500 }
    );
  }
}