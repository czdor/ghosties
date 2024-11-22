export const capitalize = (s: string) => {
  const firstLetter = s[0].toUpperCase();
  return firstLetter + s.slice(1, s.length);
};

export const sleep = (ms: number) =>
  new Promise((res, _) => setTimeout(res, ms));

export const isExpired = (timestamp: number): boolean => {
  const maxSecs = 5;
  if (Date.now() - timestamp > maxSecs * 1000) return true;

  return false;
};

export const debounce = <T extends (...args: any[]) => void>(
  func: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: NodeJS.Timeout | undefined;

  return function (...args: Parameters<T>): void {
    if (timeoutId) clearTimeout(timeoutId);

    timeoutId = setTimeout(() => {
      // @ts-ignore
      func.apply(this, args);
    }, delay);
  };
};
