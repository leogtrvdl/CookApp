import { useWindowDimensions } from 'react-native';

export type Breakpoint = 'mobile' | 'tablet' | 'desktop';

export function useBreakpoint() {
  const { width, height } = useWindowDimensions();

  const breakpoint: Breakpoint =
    width >= 1024 ? 'desktop' : width >= 640 ? 'tablet' : 'mobile';

  const isMobile = breakpoint === 'mobile';
  const isTablet = breakpoint === 'tablet';
  const isDesktop = breakpoint === 'desktop';

  const cardImageHeight = isMobile ? Math.min(width * 0.55, 240) : isTablet ? 200 : 180;

  const listColumnWidth: number | `${number}%` = isDesktop ? Math.min(width * 0.45, 520) : '100%';

  return {
    width,
    height,
    breakpoint,
    isMobile,
    isTablet,
    isDesktop,
    cardImageHeight,
    listColumnWidth,
  };
}