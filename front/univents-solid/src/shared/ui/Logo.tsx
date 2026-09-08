import type { JSX } from '@solidjs/web';

export interface LogoProps extends JSX.HTMLAttributes<HTMLDivElement> {
  variant?: 'complete' | 'icon' | 'responsive';
  theme?: 'light' | 'dark' | 'default' | 'auto';
  priority?: boolean;
  imgClassName?: string;
}

export default function Logo(props: LogoProps) {
  const variant = () => props.variant ?? 'responsive';
  const theme = () => props.theme ?? 'auto';
  const dark = () => {
    if (theme() === 'dark') return true;
    if (theme() === 'light' || theme() === 'default') return false;
    if (typeof window === 'undefined') return false;

    return (
      document.documentElement.classList.contains('dark') ||
      window.matchMedia('(prefers-color-scheme: dark)').matches
    );
  };
  const source = (logoVariant: 'complete' | 'icon') => {
    if (theme() === 'default') return `/logo-${logoVariant}-default.svg`;
    return `/logo-${logoVariant}-${dark() ? 'dark' : 'light'}.svg`;
  };
  const image = (logoVariant: 'complete' | 'icon', className: string) => (
    <img
      src={source(logoVariant)}
      alt="Univents"
      width={logoVariant === 'complete' ? 560 : 519}
      height={logoVariant === 'complete' ? 663 : 504}
      loading={props.priority ? 'eager' : 'lazy'}
      fetchpriority={props.priority ? 'high' : 'auto'}
      class={`${className} ${props.imgClassName ?? ''}`}
    />
  );

  return (
    <div
      {...props}
      class={`relative flex h-full w-full items-center justify-center ${props.class ?? ''}`}
    >
      {variant() === 'responsive' && (
        <>
          {image('complete', 'hidden md:block')}
          {image('icon', 'block md:hidden')}
        </>
      )}
      {variant() === 'complete' && image('complete', 'block')}
      {variant() === 'icon' && image('icon', 'block')}
    </div>
  );
}
