import { nextTick } from 'vue'

/**
 * useLineNavigation - Composable for standardized keyboard navigation in line item tables.
 * Handles Arrow keys for +/- values, Enter for horizontal/vertical navigation, 
 * and focus management using data-row/data-col attributes.
 */

export interface LineNavigationOptions {
  rowCount: () => number
  columns: string[]
  onLastFieldTab?: () => void
  onLastFieldEnter?: () => void
  onAddField?: () => void
  onRemoveField?: (index: number) => void
  onUpdate?: (index: number, col: string, newValue: any) => void
}

export function useLineNavigation(options: LineNavigationOptions) {
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
        
        const currentValue = Number(lineData[col] || 0)
        const newValue = isUp ? currentValue + step : Math.max(min, currentValue - step)

        if (options.onUpdate) {
          options.onUpdate(index, col, newValue)
        } else {
          // Fallback to direct mutation if no onUpdate is provided
          lineData[col] = newValue
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

    // 4. Keyboard Deletion: Delete key (if in numeric col or with modifier)
    if (e.key === 'Delete' || (e.key === 'Backspace' && e.ctrlKey)) {
      if (options.onRemoveField) {
        // Only trigger delete if it's a "whole line" operation or we are intentional
        // For numeric inputs, Delete usually deletes a character, so we check if field is empty or use Ctrl
        const isInputEmpty = String(lineData[col]).length === 0 || lineData[col] === 0
        if (e.ctrlKey || isInputEmpty) {
          e.preventDefault()
          options.onRemoveField(index)
          
          // Focus the previous line if available, otherwise next
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
      let selector = `input[data-row="${row}"]`
      if (col) {
        selector += `[data-col="${col}"]`
      }
      
      const el = document.querySelector(selector) as HTMLInputElement
      if (el) {
        el.focus()
        el.select()
      }
    })
  }

  return {
    handleLineKeyDown,
    focusLineInput
  }
}
