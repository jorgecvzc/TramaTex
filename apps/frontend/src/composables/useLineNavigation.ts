import { nextTick } from 'vue'

/**
 * useLineNavigation - Composable for standardized keyboard navigation in line item tables.
 * Handles Arrow keys for +/- values, Enter for horizontal/vertical navigation, 
 * and focus management using data-row/data-col attributes.
 */

export interface LineNavigationOptions {
  rowCount: () => number
  columns: string[]
  prefix?: string
  onLastFieldTab?: () => void
  onLastFieldEnter?: () => void
  onAddField?: () => void
  onRemoveField?: (index: number) => void
  onUpdate?: (index: number, col: string, newValue: any) => void
}

export function useLineNavigation(options: LineNavigationOptions) {
  const rowAttr = options.prefix ? `data-${options.prefix}-row` : 'data-row'
  const colAttr = options.prefix ? `data-${options.prefix}-col` : 'data-col'

  /**
   * Main keydown handler to be attached to inputs in the line items table.
   */
  function handleLineKeyDown(e: KeyboardEvent, index: number, col: string, lineData: any) {
    const rowCount = options.rowCount()
    const isLastLine = index === rowCount - 1
    const isLastCol = col === options.columns[options.columns.length - 1]

    // 1. Arrows Up/Down as +/- for numeric fields
    const numericCols = ['quantity', 'qty', 'unitPrice', 'price', 'discountPercent', 'disc', 'discount_percent', 'unit_price', 'deliveredQuantity']
    if (numericCols.includes(col)) {
      if (e.key === 'ArrowUp' || e.key === 'ArrowDown') {
        e.preventDefault()
        const isUp = e.key === 'ArrowUp'
        const step = (col.includes('quant') || col.includes('qty')) ? 1 : 1.0
        const min = (col.includes('quant') || col.includes('qty')) ? 1 : 0
        
        // Use valueAsNumber from the input target for reliability
        const target = e.target as HTMLInputElement
        const currentValue = !isNaN(target.valueAsNumber) ? target.valueAsNumber : Number(lineData[col] || 0)
        const newValue = isUp ? currentValue + step : Math.max(min, currentValue - step)

        if (options.onUpdate) {
          options.onUpdate(index, col, newValue)
        } else {
          // Fallback to direct mutation if no onUpdate is provided
          const propToMutate = lineData[col] !== undefined ? col : (Object.keys(lineData).find(k => k.toLowerCase().includes(col.toLowerCase())) || col)
          lineData[propToMutate] = newValue
        }
        return
      }
    }

    // 2. Tab key management at the very end of the table
    if (e.key === 'Tab' && !e.shiftKey && isLastLine && isLastCol) {
      if (options.onLastFieldTab) {
        e.preventDefault()
        options.onLastFieldTab()
        return
      }
    }

    // 3. Enter key for navigation
    if (e.key === 'Enter') {
      e.preventDefault()
      const colIndex = options.columns.indexOf(col)
      
      if (colIndex < options.columns.length - 1) {
        // Next column in same row
        focusLineInput(index, options.columns[colIndex + 1])
      } else {
        // Last column of the row
        if (isLastLine) {
          if (options.onLastFieldEnter) {
            options.onLastFieldEnter()
          } else if (options.onAddField) {
            options.onAddField()
          }
        } else {
          // Next row, first column
          focusLineInput(index + 1, options.columns[0])
        }
      }
      return
    }

    // 4. Keyboard Deletion: Delete key
    if (e.key === 'Delete' || (e.key === 'Backspace' && e.ctrlKey)) {
      if (options.onRemoveField) {
        const isInputEmpty = String(lineData[col]).length === 0 || lineData[col] === 0
        if (e.ctrlKey || isInputEmpty) {
          e.preventDefault()
          options.onRemoveField(index)
          
          const nextRow = index > 0 ? index - 1 : 0
          if (rowCount > 1) {
             focusLineInput(nextRow, col)
          }
        }
      }
    }

    // 5. Shortcut for Add Line: Insert key
    if (e.key === 'Insert') {
      e.preventDefault()
      if (options.onAddField) {
        options.onAddField()
      }
    }
  }

  /**
   * Helper to focus a specific input based on data attributes.
   */
  function focusLineInput(row: number, col?: string) {
    nextTick(() => {
      let selector = `[${rowAttr}="${row}"]`
      if (col) {
        selector += `[${colAttr}="${col}"]`
      }
      
      const el = document.querySelector(selector) as HTMLElement & { select?: () => void }
      if (el) {
        el.focus()
        if (typeof el.select === 'function') {
          el.select()
        }
      }
    })
  }

  return {
    handleLineKeyDown,
    focusLineInput
  }
}
