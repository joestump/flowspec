import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import styles from './index.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
  emoji: string;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Declarative Pipelines',
    emoji: '📝',
    description: (
      <>
        Define workflows as YAML, not code. Describe what should happen — steps,
        inputs, conditions, error handling — and let the engine handle execution.
      </>
    ),
  },
  {
    title: 'Temporal Native',
    emoji: '⚡',
    description: (
      <>
        Compiles to durable Temporal workflows with automatic retries, timeouts,
        and fault tolerance. Your YAML becomes production-grade orchestration.
      </>
    ),
  },
  {
    title: 'Agent Ready',
    emoji: '🤖',
    description: (
      <>
        Built for AI agent orchestration. Chain agent steps, fan out in parallel,
        route based on conditions. LLMs can generate flowspec YAML dynamically.
      </>
    ),
  },
  {
    title: 'Open Source',
    emoji: '🔓',
    description: (
      <>
        MIT licensed and community-driven. Use it standalone or embed it in your
        own agent platform. Contributions welcome.
      </>
    ),
  },
];

function Feature({title, emoji, description}: FeatureItem) {
  return (
    <div className={clsx('col col--3')}>
      <div className="text--center" style={{fontSize: '3rem', marginBottom: '1rem'}}>
        {emoji}
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            Get Started →
          </Link>
          <Link
            className="button button--outline button--secondary button--lg"
            to="/docs/examples"
            style={{marginLeft: '1rem'}}>
            See Examples
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.title}
      description="A YAML DSL for defining and executing Temporal workflows">
      <HomepageHeader />
      <main>
        <section className={styles.features}>
          <div className="container">
            <div className="row">
              {FeatureList.map((props, idx) => (
                <Feature key={idx} {...props} />
              ))}
            </div>
          </div>
        </section>
      </main>
    </Layout>
  );
}
