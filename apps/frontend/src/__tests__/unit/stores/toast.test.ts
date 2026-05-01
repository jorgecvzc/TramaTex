import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useToastStore } from '@/stores/toast';

describe('Toast Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should add a toast', () => {
    const store = useToastStore();
    const id = store.addToast('Test message', 'success');
    
    expect(store.toasts).toHaveLength(1);
    expect(store.toasts[0]).toMatchObject({
      id,
      message: 'Test message',
      type: 'success',
      duration: 5000,
    });
  });

  it('should remove a toast', () => {
    const store = useToastStore();
    const id = store.addToast('Test message');
    expect(store.toasts).toHaveLength(1);
    
    store.removeToast(id);
    expect(store.toasts).toHaveLength(0);
  });

  it('should have helper methods for each type', () => {
    const store = useToastStore();
    
    store.success('Success');
    store.error('Error');
    store.warning('Warning');
    store.info('Info');
    
    expect(store.toasts).toHaveLength(4);
    expect(store.toasts[0].type).toBe('success');
    expect(store.toasts[1].type).toBe('error');
    expect(store.toasts[2].type).toBe('warning');
    expect(store.toasts[3].type).toBe('info');
  });

  it('should auto-remove toast after duration', async () => {
    const store = useToastStore();
    store.addToast('Fast toast', 'info', 100);
    expect(store.toasts).toHaveLength(1);
    
    await new Promise(resolve => setTimeout(resolve, 150));
    expect(store.toasts).toHaveLength(0);
  });
});
