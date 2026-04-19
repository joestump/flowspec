import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'flowspec',
  tagline: 'A YAML DSL for Temporal workflows',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://joestump.github.io',
  baseUrl: '/flowspec/',

  organizationName: 'joestump',
  projectName: 'flowspec',

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl:
            'https://github.com/joestump/flowspec/tree/main/docs-site/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'flowspec',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          href: 'https://github.com/joestump/flowspec',
          label: 'GitHub',
          position: 'right',
        },
        {
          href: 'https://pkg.go.dev/github.com/joestump/flowspec',
          label: 'pkg.go.dev',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Introduction',
              to: '/docs/intro',
            },
            {
              label: 'Spec Format',
              to: '/docs/spec-format',
            },
          ],
        },
        {
          title: 'Links',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/joestump/flowspec',
            },
            {
              label: 'pkg.go.dev',
              href: 'https://pkg.go.dev/github.com/joestump/flowspec',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Joe Stump. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['yaml', 'go', 'bash'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
