export interface PrintIssuerProfile {
  displayName: string
  taxLabel: string
  taxId: string
  addressLine: string
  cityLine: string
  contactLine: string
}

const STORAGE_KEY = 'tramatex_print_issuer_profile'

const DEFAULT_PROFILE: PrintIssuerProfile = {
  displayName: 'TramaTex',
  taxLabel: 'CIF',
  taxId: '',
  addressLine: '',
  cityLine: '',
  contactLine: '',
}

function getEnvProfile(): Partial<PrintIssuerProfile> {
  const env = import.meta.env as Record<string, string | undefined>

  return {
    displayName: env.VITE_PRINT_ISSUER_NAME,
    taxLabel: env.VITE_PRINT_ISSUER_TAX_LABEL,
    taxId: env.VITE_PRINT_ISSUER_TAX_ID,
    addressLine: env.VITE_PRINT_ISSUER_ADDRESS,
    cityLine: env.VITE_PRINT_ISSUER_CITY,
    contactLine: env.VITE_PRINT_ISSUER_CONTACT,
  }
}

function getStorageProfile(): Partial<PrintIssuerProfile> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return {}

    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') return {}

    return {
      displayName: typeof parsed.displayName === 'string' ? parsed.displayName : undefined,
      taxLabel: typeof parsed.taxLabel === 'string' ? parsed.taxLabel : undefined,
      taxId: typeof parsed.taxId === 'string' ? parsed.taxId : undefined,
      addressLine: typeof parsed.addressLine === 'string' ? parsed.addressLine : undefined,
      cityLine: typeof parsed.cityLine === 'string' ? parsed.cityLine : undefined,
      contactLine: typeof parsed.contactLine === 'string' ? parsed.contactLine : undefined,
    }
  } catch {
    return {}
  }
}

export function getPrintIssuerProfile(): PrintIssuerProfile {
  return {
    ...DEFAULT_PROFILE,
    ...getEnvProfile(),
    ...getStorageProfile(),
  }
}

export function savePrintIssuerProfile(profile: Partial<PrintIssuerProfile>): PrintIssuerProfile {
  const nextProfile: PrintIssuerProfile = {
    ...getPrintIssuerProfile(),
    ...profile,
  }

  localStorage.setItem(STORAGE_KEY, JSON.stringify(nextProfile))
  return nextProfile
}

export function resetPrintIssuerProfile(): PrintIssuerProfile {
  localStorage.removeItem(STORAGE_KEY)
  return getPrintIssuerProfile()
}
