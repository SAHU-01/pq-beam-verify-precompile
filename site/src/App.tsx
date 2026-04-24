import { useState } from 'react'

// ── Inline style helpers (non-layout) ─────────────────────────────────

const heading: React.CSSProperties = {
  fontFamily: 'var(--mono)', fontSize: 13, fontWeight: 700,
  textTransform: 'uppercase', letterSpacing: '2px', marginBottom: 16,
}

const body: React.CSSProperties = {
  fontFamily: 'var(--sans)', fontSize: 14, lineHeight: 1.7, color: '#333',
}

const mono: React.CSSProperties = {
  fontFamily: 'var(--mono)', fontSize: 12, lineHeight: 1.6,
}

const expandHint: React.CSSProperties = {
  fontFamily: 'var(--mono)', fontSize: 10, color: '#999',
  marginTop: 'auto', paddingTop: 12, letterSpacing: '0.5px',
}

const toggleIcon = (open: boolean): React.CSSProperties => ({
  ...mono, fontSize: 18, color: '#999',
  transform: open ? 'rotate(45deg)' : 'none', transition: 'transform 0.2s',
})

// Detail panel sub-styles
const dHead: React.CSSProperties = {
  fontFamily: 'var(--mono)', fontSize: 18, fontWeight: 700,
  marginBottom: 24, paddingBottom: 12, borderBottom: '2px solid #000',
}
const dSub: React.CSSProperties = {
  fontFamily: 'var(--mono)', fontSize: 13, fontWeight: 700,
  textTransform: 'uppercase', letterSpacing: '1.5px', marginTop: 32, marginBottom: 12,
}
const dCode: React.CSSProperties = {
  background: 'var(--terminal-bg)', color: 'var(--terminal-text)',
  fontFamily: 'var(--mono)', fontSize: 11, padding: '16px 20px',
  lineHeight: 1.7, border: '1px solid #000', overflow: 'auto',
  whiteSpace: 'pre', marginBottom: 16,
}
const dtH: React.CSSProperties = {
  textAlign: 'left', padding: '10px 12px', borderBottom: '2px solid #000',
  fontSize: 10, textTransform: 'uppercase', letterSpacing: '1px', fontWeight: 700,
}
const dtD: React.CSSProperties = { padding: '10px 12px', borderBottom: '1px solid #ddd' }
const dtTable: React.CSSProperties = {
  width: '100%', borderCollapse: 'collapse', fontFamily: 'var(--mono)', fontSize: 12, marginBottom: 16,
}

// ── SVG Illustrations ─────────────────────────────────────────────────

function ShieldIcon() {
  return (
    <svg width="120" height="120" viewBox="0 0 120 120" fill="none" style={{ margin: '0 auto 16px', maxWidth: '100%' }}>
      <path d="M60 10 L100 30 L100 65 Q100 95 60 110 Q20 95 20 65 L20 30 Z" stroke="#000" strokeWidth="1.5" fill="none" strokeLinejoin="round" />
      <rect x="45" y="55" width="30" height="24" rx="2" stroke="#000" strokeWidth="1.5" fill="none" />
      <path d="M52 55 L52 45 Q52 38 60 38 Q68 38 68 45 L68 55" stroke="#000" strokeWidth="1.5" fill="none" />
      <circle cx="60" cy="67" r="3" fill="#000" />
      <circle cx="38" cy="42" r="2" fill="#000" opacity="0.3" />
      <circle cx="82" cy="42" r="2" fill="#000" opacity="0.3" />
      <circle cx="60" cy="22" r="2" fill="#000" opacity="0.3" />
      <line x1="38" y1="42" x2="60" y2="22" stroke="#000" strokeWidth="0.5" opacity="0.2" />
      <line x1="82" y1="42" x2="60" y2="22" stroke="#000" strokeWidth="0.5" opacity="0.2" />
    </svg>
  )
}

function LayerStackIcon() {
  return (
    <svg width="140" height="120" viewBox="0 0 140 120" fill="none" style={{ margin: '0 auto', maxWidth: '100%' }}>
      <path d="M30 85 L70 100 L110 85 L70 70 Z" stroke="#000" strokeWidth="1.5" fill="#f5f5f5" />
      <text x="70" y="88" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#333">liboqs</text>
      <path d="M30 65 L70 80 L110 65 L70 50 Z" stroke="#000" strokeWidth="1.5" fill="#eee" />
      <text x="70" y="68" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#333">CGo Bridge</text>
      <path d="M30 45 L70 60 L110 45 L70 30 Z" stroke="#000" strokeWidth="1.5" fill="#ddd" />
      <text x="70" y="48" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fontWeight="bold" fill="#000">PQ_VERIFY</text>
      <path d="M70 20 L70 10 M66 14 L70 10 L74 14" stroke="#000" strokeWidth="1.5" />
      <text x="70" y="7" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#666">Smart Contracts</text>
      <line x1="30" y1="45" x2="30" y2="85" stroke="#000" strokeWidth="0.5" strokeDasharray="2,2" />
      <line x1="110" y1="45" x2="110" y2="85" stroke="#000" strokeWidth="0.5" strokeDasharray="2,2" />
    </svg>
  )
}

function FlowDiagramIcon() {
  return (
    <svg width="100%" height="60" viewBox="0 0 500 60" fill="none" preserveAspectRatio="xMidYMid meet" style={{ marginTop: 8 }}>
      <rect x="0" y="10" width="90" height="40" stroke="#000" strokeWidth="1.5" fill="none" rx="2" />
      <text x="45" y="34" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#000">Contract</text>
      <rect x="130" y="10" width="90" height="40" stroke="#000" strokeWidth="1.5" fill="none" rx="2" />
      <text x="175" y="34" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#000">staticcall</text>
      <rect x="260" y="10" width="90" height="40" stroke="#000" strokeWidth="1.5" fill="#f0f0f0" rx="2" />
      <text x="305" y="34" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fontWeight="bold" fill="#000">PQ_VERIFY</text>
      <rect x="390" y="10" width="90" height="40" stroke="#000" strokeWidth="1.5" fill="none" rx="2" />
      <text x="435" y="30" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#000">true /</text>
      <text x="435" y="42" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#000">false</text>
      <line x1="90" y1="30" x2="130" y2="30" stroke="#000" strokeWidth="1.5" />
      <polygon points="127,26 134,30 127,34" fill="#000" />
      <line x1="220" y1="30" x2="260" y2="30" stroke="#000" strokeWidth="1.5" />
      <polygon points="257,26 264,30 257,34" fill="#000" />
      <line x1="350" y1="30" x2="390" y2="30" stroke="#000" strokeWidth="1.5" />
      <polygon points="387,26 394,30 387,34" fill="#000" />
    </svg>
  )
}

function AlgorithmIcons() {
  return (
    <svg width="160" height="100" viewBox="0 0 160 100" fill="none" style={{ margin: '0 auto 12px', maxWidth: '100%' }}>
      <text x="40" y="12" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#666">ML-DSA</text>
      {[0,1,2,3].map(r => [0,1,2,3].map(c => (
        <circle key={`${r}${c}`} cx={20 + c*14} cy={22 + r*14} r="2" fill="#000" opacity="0.6" />
      )))}
      {[0,1,2].map(r => [0,1,2].map(c => (
        <line key={`l${r}${c}`} x1={20+c*14} y1={22+r*14} x2={20+(c+1)*14} y2={22+(r+1)*14} stroke="#000" strokeWidth="0.5" opacity="0.3" />
      )))}
      <text x="120" y="12" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#666">SLH-DSA</text>
      <circle cx="120" cy="28" r="4" stroke="#000" strokeWidth="1" fill="none" />
      <line x1="120" y1="32" x2="106" y2="46" stroke="#000" strokeWidth="1" />
      <line x1="120" y1="32" x2="134" y2="46" stroke="#000" strokeWidth="1" />
      <circle cx="106" cy="48" r="3" stroke="#000" strokeWidth="1" fill="none" />
      <circle cx="134" cy="48" r="3" stroke="#000" strokeWidth="1" fill="none" />
      <line x1="106" y1="51" x2="98" y2="62" stroke="#000" strokeWidth="0.7" />
      <line x1="106" y1="51" x2="114" y2="62" stroke="#000" strokeWidth="0.7" />
      <line x1="134" y1="51" x2="126" y2="62" stroke="#000" strokeWidth="0.7" />
      <line x1="134" y1="51" x2="142" y2="62" stroke="#000" strokeWidth="0.7" />
      {[98,114,126,142].map(x => (
        <circle key={x} cx={x} cy={64} r="2.5" stroke="#000" strokeWidth="0.7" fill="none" />
      ))}
      <rect x="50" y="78" width="60" height="16" rx="2" stroke="#000" strokeWidth="1" fill="none" />
      <text x="80" y="89" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fontWeight="bold" fill="#000">NIST FIPS</text>
    </svg>
  )
}

function PuzzleIcon() {
  return (
    <svg width="100" height="80" viewBox="0 0 100 80" fill="none" style={{ margin: '0 auto 8px', maxWidth: '100%' }}>
      <rect x="10" y="10" width="30" height="25" stroke="#000" strokeWidth="1.5" fill="#f5f5f5" rx="1" />
      <text x="25" y="26" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#000">EVM</text>
      <rect x="42" y="10" width="30" height="25" stroke="#000" strokeWidth="1.5" fill="#eee" rx="1" />
      <text x="57" y="26" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#000">Sol</text>
      <rect x="10" y="37" width="30" height="25" stroke="#000" strokeWidth="1.5" fill="#eee" rx="1" />
      <text x="25" y="53" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#000">Go</text>
      <rect x="42" y="37" width="30" height="25" stroke="#000" strokeWidth="1.5" fill="#f5f5f5" rx="1" />
      <text x="57" y="53" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#000">CGo</text>
      <text x="82" y="42" fontFamily="var(--mono)" fontSize="16" fill="#000">+</text>
    </svg>
  )
}

function BeamArchDiagram() {
  return (
    <svg width="100%" viewBox="0 0 520 260" fill="none" preserveAspectRatio="xMidYMid meet" style={{ maxWidth: 520, margin: '0 auto', display: 'block' }}>
      {/* Avalanche Primary Network */}
      <rect x="10" y="10" width="500" height="50" stroke="#000" strokeWidth="1.5" fill="#f5f5f5" rx="2" />
      <text x="260" y="30" textAnchor="middle" fontFamily="var(--mono)" fontSize="10" fontWeight="bold" fill="#000">AVALANCHE PRIMARY NETWORK</text>
      <text x="260" y="46" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#666">P-Chain / X-Chain / C-Chain</text>
      {/* Arrow */}
      <line x1="260" y1="60" x2="260" y2="80" stroke="#000" strokeWidth="1.5" />
      <polygon points="256,77 260,83 264,77" fill="#000" />
      {/* Beam Subnet */}
      <rect x="60" y="85" width="400" height="50" stroke="#000" strokeWidth="1.5" fill="#eee" rx="2" />
      <text x="260" y="105" textAnchor="middle" fontFamily="var(--mono)" fontSize="10" fontWeight="bold" fill="#000">BEAM SUBNET (Subnet-EVM v0.8.0)</text>
      <text x="260" y="121" textAnchor="middle" fontFamily="var(--mono)" fontSize="8" fill="#666">Gaming / NFTs / DeFi — Chain ID 13337 — 4,500 TPS</text>
      {/* Arrow */}
      <line x1="260" y1="135" x2="260" y2="155" stroke="#000" strokeWidth="1.5" />
      <polygon points="256,152 260,158 264,152" fill="#000" />
      {/* PQ_VERIFY Precompile */}
      <rect x="120" y="160" width="280" height="40" stroke="#000" strokeWidth="2" fill="#ddd" rx="2" />
      <text x="260" y="184" textAnchor="middle" fontFamily="var(--mono)" fontSize="10" fontWeight="bold" fill="#000">PQ_VERIFY PRECOMPILE @ 0x0300...0000</text>
      {/* Arrow */}
      <line x1="260" y1="200" x2="260" y2="218" stroke="#000" strokeWidth="1.5" />
      <polygon points="256,215 260,221 264,215" fill="#000" />
      {/* liboqs */}
      <rect x="160" y="222" width="200" height="30" stroke="#000" strokeWidth="1.5" fill="#f5f5f5" rx="2" />
      <text x="260" y="241" textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#333">liboqs 0.15 (CGo) — ML-DSA + SLH-DSA</text>
      {/* Side labels */}
      <text x="30" y="180" fontFamily="var(--mono)" fontSize="7" fill="#666" transform="rotate(-90 30 180)">NATIVE</text>
      <text x="490" y="180" fontFamily="var(--mono)" fontSize="7" fill="#666" transform="rotate(90 490 180)">GAS-METERED</text>
    </svg>
  )
}

function TransactionFlowDiagram() {
  const bW = 200, bH = 34, cx = 280

  const Box = ({ y, label, fill }: { y: number; label: string; fill?: string }) => (
    <g>
      <rect x={cx - bW / 2} y={y} width={bW} height={bH} stroke="#000" strokeWidth="1.5" fill={fill || 'none'} />
      <text x={cx} y={y + bH / 2 + 4} textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fill="#000">{label}</text>
    </g>
  )
  const Arr = ({ y1, y2 }: { y1: number; y2: number }) => (
    <g>
      <line x1={cx} y1={y1} x2={cx} y2={y2} stroke="#000" strokeWidth="1.5" />
      <polygon points={`${cx-4},${y2-4} ${cx},${y2} ${cx+4},${y2-4}`} fill="#000" />
    </g>
  )

  return (
    <svg width="100%" viewBox="0 0 560 440" fill="none" preserveAspectRatio="xMidYMid meet" className="flow-svg">
      <text x="20" y="18" fontFamily="var(--mono)" fontSize="10" fontWeight="bold">EIP-2718 TYPE 0x50 TRANSACTION FLOW</text>
      <rect x={cx - bW / 2 - 15} y="35" width={bW + 30} height="75" stroke="#000" strokeWidth="1" fill="none" strokeDasharray="4,3" />
      <text x={cx} y="48" textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fontWeight="bold" fill="#666">SMART CONTRACT LAYER</text>
      <Box y={55} label="SMART CONTRACT" fill="#f0f0f0" />
      <text x={cx} y={55 + bH + 12} textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#666">Constructs Type 0x50 tx with PQ data</text>
      <Arr y1={105} y2={130} />
      <Box y={130} label="SUBMISSION TO NETWORK" />
      <text x={20} y={145} fontFamily="var(--mono)" fontSize="7" fill="#666">0x50 || RLP([chainId,</text>
      <text x={20} y={156} fontFamily="var(--mono)" fontSize="7" fill="#666">nonce, ...data, pqAlg,</text>
      <text x={20} y={167} fontFamily="var(--mono)" fontSize="7" fill="#666">pqPub, pqSig])</text>
      <Arr y1={164} y2={190} />
      <Box y={190} label="EXECUTION IN EVM" />
      <Arr y1={224} y2={250} />
      <Box y={250} label="staticcall → PQ_VERIFY" fill="#f0f0f0" />
      <text x={20} y={265} fontFamily="var(--mono)" fontSize="7" fill="#666">Precompile at</text>
      <text x={20} y={276} fontFamily="var(--mono)" fontSize="7" fill="#666">0x0300...0000</text>
      <Arr y1={284} y2={310} />
      <Box y={310} label="C BRIDGE (CGo)" />
      {/* Input boxes */}
      <text x={cx + bW / 2 + 14} y={306} fontFamily="var(--mono)" fontSize="7" fill="#666">CALL INPUTS</text>
      {['MSG', 'PubKey', 'Sig'].map((l, i) => (
        <g key={l}>
          <rect x={cx + bW / 2 + 10 + i * 56} y={312} width={50} height={22} stroke="#000" strokeWidth="1" fill="none" />
          <text x={cx + bW / 2 + 35 + i * 56} y={326} textAnchor="middle" fontFamily="var(--mono)" fontSize="7" fill="#333">{l}</text>
        </g>
      ))}
      <line x1={cx + bW / 2} y1={323} x2={cx + bW / 2 + 10} y2={323} stroke="#000" strokeWidth="0.7" />
      <Arr y1={344} y2={370} />
      <Box y={370} label="liboqs VERIFICATION" fill="#f0f0f0" />
      <Arr y1={404} y2={420} />
      <text x={cx} y={435} textAnchor="middle" fontFamily="var(--mono)" fontSize="9" fontWeight="bold">abi.encode(bool valid)</text>
    </svg>
  )
}

// ── Expandable Detail Panels ──────────────────────────────────────────

function PerformanceDetail() {
  return (
    <div className="detail-panel">
      <div style={dHead}>Performance Research Log</div>
      <p style={body}>All benchmarks use real cryptographic operations via liboqs through CGo. No simulations, no mocks.</p>

      <div style={dSub}>Test Environment</div>
      <div style={dCode}>{`Platform:     Apple M1 Pro (ARM64)
OS:           macOS (Darwin)
Go:           1.23+
CGO_ENABLED:  1
liboqs:       0.15.0 (via Homebrew)
OpenSSL:      3.x (required by liboqs)
Iterations:   1,000 per operation
Data:         Ephemeral keys generated per run`}</div>

      <div style={dSub}>Benchmark Results</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr>
            <th style={dtH}>Operation</th><th style={dtH}>Avg Time</th><th style={dtH}>Ops/sec</th><th style={dtH}>Memory</th><th style={dtH}>vs ecrecover</th>
          </tr></thead>
          <tbody>
            {[
              { op: 'ML-DSA-65 verify', time: '105 us', ops: '9,504', mem: '0 B/op', vs: '4.2x' },
              { op: 'ML-DSA-65 sign', time: '497 us', ops: '2,013', mem: '3,488 B/op', vs: '19.9x' },
              { op: 'SLH-DSA-128s verify', time: '432 us', ops: '2,316', mem: '0 B/op', vs: '17.3x' },
              { op: 'SLH-DSA-128s sign', time: '439 ms', ops: '2.3', mem: '8,224 B/op', vs: '17,560x' },
              { op: 'SHA256+Keccak (baseline)', time: '0.5 us', ops: '2,051,636', mem: '32 B/op', vs: '0.02x' },
            ].map((r, i) => (
              <tr key={i} style={i === 4 ? { background: '#f5f5f5' } : {}}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.op}</td>
                <td style={dtD}>{r.time}</td><td style={dtD}>{r.ops}</td><td style={dtD}>{r.mem}</td><td style={dtD}>{r.vs}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div style={dSub}>Go Benchmark Output</div>
      <div style={dCode}>{`$ CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=. -benchmem

goos: darwin
goarch: arm64
pkg: pq-beam-verify-precompile/cmd/benchmark

BenchmarkMLDSA65Verify-10       9504    105234 ns/op      0 B/op    0 allocs/op
BenchmarkMLDSA65Sign-10         2013    497012 ns/op   3488 B/op   12 allocs/op
BenchmarkSLHDSA128sVerify-10    2316    432187 ns/op      0 B/op    0 allocs/op
BenchmarkSLHDSA128sSign-10         2 439012345 ns/op   8224 B/op   34 allocs/op

PASS
ok  pq-beam-verify-precompile/cmd/benchmark  12.847s`}</div>

      <div style={dSub}>Gas Calculation Methodology</div>
      <div style={dCode}>{`ecrecover baseline:     ~25 us = 3,000 gas (Ethereum standard)

ML-DSA-65:
  verify time:          105 us
  ratio vs ecrecover:   105 / 25 = 4.2x
  raw gas:              4.2 x 3,000 = 12,600
  with 10x safety:      126,000 → rounded to 130,000 gas
  + base overhead:      3,600 gas
  total:                133,600 gas

SLH-DSA-128s:
  verify time:          432 us
  ratio vs ecrecover:   432 / 25 = 17.3x
  raw gas:              17.3 x 3,000 = 51,900
  with 10x safety:      519,000 → rounded to 520,000 gas
  + base overhead:      3,600 gas
  total:                523,600 gas`}</div>

      <div style={dSub}>Why 10x Safety Margin?</div>
      <div className="detail-grid-2">
        {[
          { title: 'Hardware variance', desc: 'Validators run on different CPUs — ARM vs x86, different generations, different cache sizes.' },
          { title: 'liboqs build variance', desc: 'Compilation flags, optimization levels, and platform-specific assembly affect timing.' },
          { title: 'DoS protection', desc: 'Underpriced operations can be exploited to stall validators. Higher gas prevents abuse.' },
          { title: 'Adjustable per-chain', desc: 'Can be reduced via genesis gasOverrides as hardware improves and benchmarks stabilize across validators.' },
        ].map((item, i) => (
          <div key={i} style={{ border: '1px solid #ddd', padding: 16 }}>
            <div style={{ ...mono, fontWeight: 700, marginBottom: 4 }}>{item.title}</div>
            <div style={{ fontSize: 13, color: '#555' }}>{item.desc}</div>
          </div>
        ))}
      </div>

      <div style={dSub}>Key & Signature Sizes</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr>
            <th style={dtH}>Algorithm</th><th style={dtH}>Standard</th><th style={dtH}>Pub Key</th><th style={dtH}>Secret Key</th><th style={dtH}>Signature</th><th style={dtH}>Security</th>
          </tr></thead>
          <tbody>
            {[
              { alg: 'ML-DSA-65', std: 'FIPS 204', pk: '1,952 B', sk: '4,032 B', sig: '3,309 B', sec: 'Level 3 / AES-192' },
              { alg: 'SLH-DSA-128s', std: 'FIPS 205', pk: '32 B', sk: '64 B', sig: '7,856 B', sec: 'Level 1 / AES-128' },
              { alg: 'ECDSA (current)', std: '--', pk: '33 B', sk: '32 B', sig: '65 B', sec: null },
            ].map((r, i) => (
              <tr key={i}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.alg}</td><td style={dtD}>{r.std}</td>
                <td style={dtD}>{r.pk}</td><td style={dtD}>{r.sk}</td><td style={dtD}>{r.sig}</td>
                <td style={dtD}>{r.sec === null
                  ? <span style={{ ...mono, fontSize: 10, padding: '2px 8px', border: '1px solid var(--red)', color: 'var(--red)', fontWeight: 700 }}>VULNERABLE</span>
                  : r.sec}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div style={dSub}>How to Reproduce</div>
      <div style={dCode}>{`# Install liboqs
brew install liboqs openssl@3

# Clone
git clone https://github.com/SAHU-01/pq-beam-verify-precompile.git
cd pq-beam-verify-precompile

# Run benchmarks
CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=. -benchmem -count=3

# Generate benchmarks.json (structured output)
CGO_ENABLED=1 go run ./cmd/benchmark/`}</div>
    </div>
  )
}

function OpenSourceDetail() {
  return (
    <div className="detail-panel">
      <div style={dHead}>Repository Guide</div>
      <p style={body}>The entire project is open source under the MIT license. Here's how to navigate, build, and run everything.</p>

      <div style={dSub}>Repository Structure</div>
      <div style={dCode}>{`pq-beam-verify-precompile/
|-- cmd/benchmark/         # Benchmark CLI + Go test benchmarks
|-- pkg/pqcrypto/          # CGo bindings to liboqs (keygen, sign, verify)
|-- pkg/pqverify/          # EVM precompile implementation
|-- contracts/             # IPQVerify, PQAccount (ERC-4337), PQKeyRotation
|-- test/                  # End-to-end tests (7 tests)
|-- scripts/               # deploy_local.sh, demo_onchain.sh, genesis.json
|-- docs/                  # TECHNICAL_SPEC.md, LOCAL_DEPLOYMENT.md
|-- site/                  # This landing page (React + Vite)`}</div>

      <div style={dSub}>Key Files Explained</div>
      <div className="detail-grid-4">
        {[
          { file: 'pkg/pqcrypto/pqcrypto.go', desc: 'The core. CGo bindings that call liboqs C functions for key generation, signing, and verification.' },
          { file: 'pkg/pqverify/precompile.go', desc: 'The EVM precompile. Handles ABI decoding, gas metering, and dispatching to pqcrypto.' },
          { file: 'contracts/PQAccount.sol', desc: 'Smart account controlled by a PQ public key. Validates signatures via the precompile with nonce replay protection.' },
          { file: 'scripts/genesis.json', desc: 'Subnet genesis config. Activates PQ_VERIFY at block 0, sets gas limits, chain ID 13337.' },
        ].map((item, i) => (
          <div key={i} style={{ border: '1px solid #000', padding: 16, marginTop: i >= 2 ? -1 : 0, marginLeft: i % 2 === 1 ? -1 : 0 }}>
            <div style={{ ...mono, fontWeight: 700, fontSize: 11, marginBottom: 6 }}>{item.file}</div>
            <div style={{ fontSize: 13, color: '#555', lineHeight: 1.6 }}>{item.desc}</div>
          </div>
        ))}
      </div>

      <div style={dSub}>Quick Start</div>
      <div style={dCode}>{`# Prerequisites
brew install go liboqs openssl@3

# Clone
git clone https://github.com/SAHU-01/pq-beam-verify-precompile.git
cd pq-beam-verify-precompile

# Run all tests
CGO_ENABLED=1 go test ./... -v

# Run benchmarks
CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=. -benchmem

# Deploy local subnet (requires AvalancheGo + Subnet-EVM)
./scripts/deploy_local.sh

# Run on-chain demo
RPC_URL=http://127.0.0.1:9650/ext/bc/<subnet-id>/rpc \\
  ./scripts/demo_onchain.sh`}</div>

      <div style={dSub}>Build Requirements</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr><th style={dtH}>Dependency</th><th style={dtH}>Version</th><th style={dtH}>Why</th></tr></thead>
          <tbody>
            {[
              { dep: 'Go', ver: '1.23+', why: 'Build language' },
              { dep: 'liboqs', ver: '0.15.0+', why: 'PQ cryptographic primitives' },
              { dep: 'OpenSSL', ver: '3.x', why: 'Required by liboqs for symmetric crypto' },
              { dep: 'CGO_ENABLED=1', ver: '--', why: 'Mandatory for CGo bridge to C' },
              { dep: 'AvalancheGo', ver: 'v1.14.1-antithesis', why: 'Local subnet deployment' },
              { dep: 'Subnet-EVM', ver: 'v0.8.0', why: 'Modified fork with PQ_VERIFY' },
            ].map((r, i) => (
              <tr key={i}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.dep}</td><td style={dtD}>{r.ver}</td><td style={dtD}>{r.why}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function PhaseDetail({ phase }: { phase: number }) {
  if (phase === 1) return (
    <div className="detail-panel">
      <div style={dHead}>Phase 1: Precompile — In Progress</div>
      <p style={body}>The cryptographic foundation is built. MVP is live with 36+ tests and on-chain proof. Remaining: Beam testnet fork deployment, gas calibration on validator hardware, and integration testing.</p>
      <div style={dSub}>Milestones Delivered</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr><th style={dtH}>Milestone</th><th style={dtH}>Status</th><th style={dtH}>Details</th></tr></thead>
          <tbody>
            {[
              { ms: 'PQ_VERIFY precompile', detail: 'Registered at 0x0300...0000. ML-DSA-65 + SLH-DSA-128s via CGo/liboqs.', done: true },
              { ms: 'ERC-4337 smart account', detail: 'PQAccount with validateUserOp(), batch execution, EntryPoint integration.', done: true },
              { ms: 'Key rotation contract', detail: 'PQKeyRotation.sol — ECDSA→PQ migration, PQ→PQ rotation with 24h timelock.', done: true },
              { ms: 'Fuzz-tested ABI decoder', detail: 'FuzzDecodeInput + FuzzPrecompileRun — found and fixed overflow bug.', done: true },
              { ms: 'On-chain proof (local)', detail: 'Valid + tampered sig verified on local subnet. 36+ tests passing.', done: true },
              { ms: 'Beam testnet fork deploy', detail: 'Deploy precompile on Beam testnet fork with real validators.', done: false },
              { ms: 'Gas calibration', detail: 'Benchmark on Beam validator hardware. Calibrate safety margins.', done: false },
              { ms: 'Technical specification', detail: 'Published spec for precompile, account format, gas schedule.', done: false },
            ].map((r, i) => (
              <tr key={i}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.ms}</td>
                <td style={dtD}><span style={{ color: r.done ? 'var(--green)' : '#666', fontWeight: 700 }}>{r.done ? '+' : '-'}</span> {r.done ? 'Done' : 'Remaining'}</td>
                <td style={{ ...dtD, fontSize: 11 }}>{r.detail}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div style={dSub}>Key Decisions</div>
      <div className="detail-grid-2">
        {[
          { q: 'Why a precompile, not a Solidity library?', a: 'PQ verification requires C-level crypto (liboqs). Pure Solidity cannot call C libraries. A precompile provides native performance.' },
          { q: 'Why dual algorithms?', a: 'ML-DSA-65 is fast but newer. SLH-DSA-128s uses only hash functions — if lattice math breaks, the fallback is ready.' },
          { q: 'Why 0x0300...0000?', a: 'Follows Subnet-EVM custom precompile range. Avoids conflict with Ethereum native precompiles (0x01-0x09).' },
          { q: 'Why staticcall only?', a: 'Verification is pure — no state changes. staticcall prevents reentrancy and storage corruption.' },
          { q: 'Why ERC-4337?', a: 'Account abstraction is the standard for smart accounts. Implementing IAccount.validateUserOp() makes PQ accounts compatible with bundlers, paymasters, and the existing AA ecosystem.' },
          { q: 'Why a key rotation timelock?', a: 'If a PQ key is compromised, the 24h timelock gives the owner time to cancel the rotation via their ECDSA fallback (revokeKey). Defense in depth.' },
        ].map((item, i) => (
          <div key={i} style={{ border: '1px solid #ddd', padding: 16 }}>
            <div style={{ ...mono, fontWeight: 700, fontSize: 11, marginBottom: 6 }}>{item.q}</div>
            <div style={{ fontSize: 13, color: '#555', lineHeight: 1.6 }}>{item.a}</div>
          </div>
        ))}
      </div>
      <div style={{ ...mono, fontSize: 12, color: '#555', marginTop: 24 }}>Timeline: May — July 2026</div>
    </div>
  )

  if (phase === 2) return (
    <div className="detail-panel">
      <div style={dHead}>Phase 2: SDK + Native Support — Next</div>
      <p style={body}>Making PQ signatures usable by application developers. A TypeScript SDK will handle key management and signing.</p>
      <div style={dSub}>Planned Milestones</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr><th style={dtH}>Milestone</th><th style={dtH}>Status</th><th style={dtH}>Details</th></tr></thead>
          <tbody>
            {[
              { ms: 'TypeScript SDK', detail: 'npm package for PQ key generation, signing, and tx construction.' },
              { ms: 'Key management', detail: 'Secure storage, derivation, backup/recovery for PQ keypairs.' },
              { ms: 'Type 0x50 in validators', detail: 'Validators natively parse and verify PQ transactions.' },
              { ms: 'Auto PQ key creation', detail: 'Generate PQ keypair alongside ECDSA. Dual-signature migration.' },
            ].map((r, i) => (
              <tr key={i}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.ms}</td>
                <td style={dtD}><span style={{ color: '#666' }}>-</span> Planned</td>
                <td style={{ ...dtD, fontSize: 11 }}>{r.detail}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div style={dSub}>Open Considerations</div>
      <div className="detail-grid-2">
        {[
          { title: 'WASM vs Native', desc: 'liboqs compiled to WASM for browsers? Or native N-API bindings? WASM is portable but slower.' },
          { title: 'Key derivation', desc: 'PQ keys are much larger. BIP-39 mnemonics may not map cleanly. May need PQ-specific derivation.' },
          { title: 'Wallet UX', desc: 'PQ public keys are 1,952 bytes vs 33 bytes. Wallet UI must handle display and QR codes for large keys.' },
          { title: 'Migration path', desc: 'Dual-signature period where both ECDSA and PQ are valid, then gradual ECDSA sunset.' },
        ].map((item, i) => (
          <div key={i} style={{ border: '1px solid #ddd', padding: 16 }}>
            <div style={{ ...mono, fontWeight: 700, fontSize: 11, marginBottom: 6 }}>{item.title}</div>
            <div style={{ fontSize: 13, color: '#555', lineHeight: 1.6 }}>{item.desc}</div>
          </div>
        ))}
      </div>
      <div style={{ ...mono, fontSize: 12, color: '#555', marginTop: 24 }}>Timeline: August — October 2026</div>
    </div>
  )

  return (
    <div className="detail-panel">
      <div style={dHead}>Phase 3: Audit + Mainnet — Planned</div>
      <p style={body}>Production readiness: security audit, validator benchmarks, and mainnet deployment.</p>
      <div style={dSub}>Planned Milestones</div>
      <div style={{ overflowX: 'auto' }}>
        <table style={dtTable}>
          <thead><tr><th style={dtH}>Milestone</th><th style={dtH}>Target</th><th style={dtH}>Details</th></tr></thead>
          <tbody>
            {[
              { ms: 'Security audit', target: 'M3', detail: 'CGo boundary, ABI decoder, gas accounting. Side-channel analysis.' },
              { ms: 'Validator benchmarks', target: 'M3', detail: 'Benchmarks on actual validator hardware (x86 + ARM).' },
              { ms: 'Mainnet deployment', target: 'M4', detail: 'Deploy via coordinated Beam network upgrade.' },
              { ms: 'Migration toolkit', target: 'M4', detail: 'CLI tools for migrating existing ECDSA accounts to PQ.' },
              { ms: 'Documentation', target: 'M4', detail: 'Developer docs, integration guides, security model docs.' },
            ].map((r, i) => (
              <tr key={i}>
                <td style={{ ...dtD, fontWeight: 700 }}>{r.ms}</td><td style={dtD}>{r.target}</td>
                <td style={{ ...dtD, fontSize: 11 }}>{r.detail}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div style={dSub}>Audit Scope</div>
      <div className="detail-grid-2">
        {[
          { title: 'CGo boundary', desc: 'Memory management across Go/C. Buffer allocation, deallocation, pointer safety.' },
          { title: 'ABI decoder fuzzing', desc: 'Malformed inputs, truncated data, oversized fields. No panics allowed.' },
          { title: 'Gas cost validation', desc: 'Verify costs on production hardware. 10x safety margin may be adjustable.' },
          { title: 'Smart account review', desc: 'PQAccount authorization flow, nonce handling, replay protection.' },
        ].map((item, i) => (
          <div key={i} style={{ border: '1px solid #ddd', padding: 16 }}>
            <div style={{ ...mono, fontWeight: 700, fontSize: 11, marginBottom: 6 }}>{item.title}</div>
            <div style={{ fontSize: 13, color: '#555', lineHeight: 1.6 }}>{item.desc}</div>
          </div>
        ))}
      </div>
      <div style={{ ...mono, fontSize: 12, color: '#555', marginTop: 24 }}>Timeline: November 2026 — January 2027</div>
    </div>
  )
}

// ── Main Page ─────────────────────────────────────────────────────────

function App() {
  const [showRun, setShowRun] = useState(false)
  const [expanded, setExpanded] = useState<string | null>(null)

  const toggle = (section: string) => setExpanded(expanded === section ? null : section)

  return (
    <div>
      {/* ── Title ──────────────────────────────────────────────── */}
      <div style={{ borderBottom: '2px solid #000', paddingBottom: 24 }}>
        <div style={{ fontFamily: 'var(--mono)', fontSize: 11, letterSpacing: '3px', textTransform: 'uppercase', color: '#666', marginBottom: 8 }}>
          Beam Network / Post-Quantum Cryptography
        </div>
        <h1 style={{ fontFamily: 'var(--mono)', fontSize: 36, fontWeight: 700, letterSpacing: '-1px', margin: 0 }}>
          PQ_VERIFY
        </h1>
        <p style={{ fontFamily: 'var(--sans)', fontSize: 16, color: '#333', marginTop: 8, maxWidth: 600, lineHeight: 1.6 }}>
          A precompile that lets smart contracts verify post-quantum
          signatures directly on-chain. Built for Beam's Subnet-EVM.
        </p>
      </div>

      {/* ══ ROW 1: Core intro ═════════════════════════════════════ */}
      <div className="bento-3">
        <div className="cell">
          <div style={heading}>Simple and Secure</div>
          <p style={body}>
            Quantum computers will break today's blockchain
            signatures. PQ_VERIFY adds quantum-resistant
            verification as a native EVM operation.
          </p>
          <p style={{ ...body, marginTop: 12 }}>
            Any smart contract can call it. No libraries to import,
            no off-chain services. One <code style={mono}>staticcall</code>.
          </p>
          <ShieldIcon />
        </div>

        <div className="terminal-cell" style={{ borderTop: 'none', borderBottom: 'none' }}>
          <div style={{ color: '#888', marginBottom: 12, fontSize: 11 }}>// Verify a post-quantum signature from Solidity</div>
          <div style={{ marginBottom: 16 }}>
            <span style={{ color: '#888' }}>{'>'}</span>{' '}
            <span style={{ color: 'var(--terminal-green)' }}>address</span> PQ_VERIFY =<br />
            {'    '}0x0300...0000;
          </div>
          <div style={{ color: '#ccc', lineHeight: 2 }}>
            <span style={{ color: '#888' }}>(</span>bool success, bytes memory out<span style={{ color: '#888' }}>)</span> =<br />
            {'  '}PQ_VERIFY.<span style={{ color: 'var(--terminal-green)' }}>staticcall</span><span style={{ color: '#888' }}>(</span><br />
            {'    '}abi.encode<span style={{ color: '#888' }}>(</span><br />
            {'      '}pubkey,{'     '}<span style={{ color: '#888' }}>// bytes</span><br />
            {'      '}signature,{'  '}<span style={{ color: '#888' }}>// bytes</span><br />
            {'      '}message,{'    '}<span style={{ color: '#888' }}>// bytes</span><br />
            {'      '}algorithm{'   '}<span style={{ color: '#888' }}>// 0 = ML-DSA-65</span><br />
            {'    '}<span style={{ color: '#888' }}>)</span><br />
            {'  '}<span style={{ color: '#888' }}>)</span>;<br /><br />
            bool <span style={{ color: 'var(--terminal-green)' }}>valid</span> = abi.decode<span style={{ color: '#888' }}>(</span>out, <span style={{ color: '#888' }}>(</span>bool<span style={{ color: '#888' }}>)</span><span style={{ color: '#888' }}>)</span>;
          </div>
        </div>

        <div className="cell">
          <div style={heading}>Two Algorithms</div>
          <p style={body}>Both are NIST-standardized (2024) and production-ready.</p>
          <AlgorithmIcons />
          <div style={{ marginTop: 12 }}>
            <div style={{ ...mono, fontWeight: 700, marginBottom: 4 }}>ML-DSA-65</div>
            <div style={{ ...mono, color: '#666', fontSize: 11 }}>Lattice-based. Fast. FIPS 204.</div>
            <div style={{ ...mono, fontWeight: 700, marginTop: 12, marginBottom: 4 }}>SLH-DSA-128s</div>
            <div style={{ ...mono, color: '#666', fontSize: 11 }}>Hash-based fallback. Conservative. FIPS 205.</div>
          </div>
        </div>
      </div>

      {/* ══ WHY BEAM ══════════════════════════════════════════════ */}
      <div className="bento-2">
        <div className="cell" style={{ borderRight: 'none' }}>
          <div style={heading}>Why Beam?</div>
          <p style={body}>
            Beam is a Subnet-EVM chain on Avalanche, purpose-built for gaming.
            Its assets — in-game items, NFTs, DeFi positions, user accounts — are
            all secured by ECDSA signatures that quantum computers will break.
          </p>
          <p style={{ ...body, marginTop: 12 }}>
            As a Subnet-EVM chain, Beam can deploy custom precompiles at the
            VM level — something <strong>impossible on Ethereum mainnet</strong> without
            a network-wide hard fork. This makes Beam the ideal proving
            ground for post-quantum infrastructure.
          </p>
          <div style={{ ...mono, fontSize: 11, marginTop: 16, padding: '12px 16px', border: '1px solid #000', background: '#f8f8f8' }}>
            "Beam is positioned to lead — as a Subnet-EVM chain, deploying a
            custom precompile is significantly simpler than on Ethereum mainnet."
          </div>
        </div>
        <div className="cell">
          <div style={heading}>What's at Stake</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 0 }}>
            {[
              { asset: 'Gaming Assets', risk: 'In-game items, skins, currencies — all controlled by ECDSA keys vulnerable to Shor\'s algorithm.' },
              { asset: 'NFT Ownership', risk: 'Provenance and ownership proofs rely on ECDSA. A quantum attacker could forge transfers.' },
              { asset: 'DeFi Positions', risk: 'Liquidity pools, staking, and lending positions secured by breakable signatures.' },
              { asset: 'Validator Keys', risk: 'Validator signing keys are ECDSA. Compromised validators threaten network consensus.' },
              { asset: 'User Accounts', risk: 'Every wallet, every EOA — derived from ECDSA key pairs that quantum computers will crack.' },
            ].map((item, i) => (
              <div key={i} style={{ padding: '10px 0', borderBottom: i < 4 ? '1px solid #ddd' : 'none' }}>
                <div style={{ ...mono, fontWeight: 700, fontSize: 11 }}>{item.asset}</div>
                <div style={{ fontSize: 13, color: '#555', lineHeight: 1.5, marginTop: 2 }}>{item.risk}</div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* ══ BEAM ARCHITECTURE ════════════════════════════════════ */}
      <div style={{ borderBottom: '1px solid #000', padding: '28px' }}>
        <div style={heading}>Beam Architecture + PQ_VERIFY Integration</div>
        <p style={{ ...body, marginBottom: 20 }}>
          PQ_VERIFY lives inside the Subnet-EVM layer — the same level as native EVM opcodes
          like ecrecover. It's not a smart contract; it's a VM-level extension.
        </p>
        <BeamArchDiagram />
        <div className="bento-3" style={{ marginTop: 24, borderBottom: 'none' }}>
          {[
            { title: 'Subnet Advantage', desc: 'Beam controls its own VM. No need for Ethereum governance or hard fork coordination. Deploy when ready.' },
            { title: 'Backward Compatible', desc: 'ECDSA transactions continue unchanged. Existing contracts need zero modifications. PQ is opt-in.' },
            { title: 'SDK Transparent', desc: 'Beam SDK social login flow is unchanged. PQ signing is handled at the infrastructure level, invisible to users.' },
          ].map((item, i) => (
            <div key={i} className="cell">
              <div style={{ ...mono, fontWeight: 700, fontSize: 11, marginBottom: 8 }}>{item.title}</div>
              <div style={{ fontSize: 13, color: '#555', lineHeight: 1.6 }}>{item.desc}</div>
            </div>
          ))}
        </div>
      </div>

      {/* ══ HOW INTEGRATION WORKS ════════════════════════════════ */}
      <div className="bento-3">
        <div className="cell">
          <div style={heading}>For Game Studios</div>
          <p style={body}>
            Call PQ_VERIFY from any Solidity contract — the same way you'd call ecrecover.
            Use PQAccount.sol as a reference for PQ-secured smart accounts.
          </p>
          <div style={{ ...mono, fontSize: 11, marginTop: 12, lineHeight: 1.8, color: '#555' }}>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>1. Import IPQVerify.sol interface</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>2. staticcall to 0x0300...0000</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>3. Pass pubkey + sig + msg + algorithm</div>
            <div style={{ padding: '4px 0' }}>4. Get back true/false</div>
          </div>
        </div>
        <div className="cell" style={{ borderLeft: 'none', borderRight: 'none' }}>
          <div style={heading}>For Beam SDK Users</div>
          <p style={body}>
            Nothing changes on the surface. The Beam SDK will handle PQ key generation
            and signing behind the scenes. Social login, session keys, and account
            abstraction all work the same way.
          </p>
          <div style={{ ...mono, fontSize: 11, marginTop: 12, lineHeight: 1.8, color: '#555' }}>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>SDK generates PQ keypair alongside ECDSA</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>Transactions auto-signed with PQ key</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>Type 0x50 envelope handled transparently</div>
            <div style={{ padding: '4px 0' }}>Gradual migration — dual-sig period</div>
          </div>
        </div>
        <div className="cell">
          <div style={heading}>For Validators</div>
          <p style={body}>
            Validators need a Subnet-EVM build with PQ_VERIFY registered and liboqs
            linked. Phase 2 adds native Type 0x50 transaction validation at the
            consensus layer.
          </p>
          <div style={{ ...mono, fontSize: 11, marginTop: 12, lineHeight: 1.8, color: '#555' }}>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>Subnet-EVM fork with PQ_VERIFY</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>CGO_ENABLED=1 build required</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>liboqs statically linked for production</div>
            <div style={{ padding: '4px 0' }}>Gas costs validated per hardware</div>
          </div>
        </div>
      </div>

      {/* ══ SCOPE OF SUCCESS ═════════════════════════════════════ */}
      <div className="bento-2">
        <div className="cell" style={{ borderRight: 'none' }}>
          <div style={heading}>Scope of Success</div>
          <p style={body}>
            If this project succeeds, Beam becomes the <strong>first Subnet-EVM chain</strong> with
            native post-quantum signature verification — demonstrating technical leadership
            in quantum readiness years ahead of Ethereum mainnet.
          </p>
          <div style={{ marginTop: 16 }}>
            {[
              { label: 'First-mover', desc: 'No other Subnet-EVM chain has PQ verification. Beam sets the standard.' },
              { label: 'Asset protection', desc: 'Billions in gaming assets, NFTs, and DeFi positions quantum-proofed before the threat materializes.' },
              { label: 'No disruption', desc: 'Backward compatible. Existing users, contracts, and tooling continue unchanged.' },
              { label: 'Gradual migration', desc: 'Opt-in adoption. Dual-signature period. No forced migration.' },
              { label: 'Throughput safe', desc: 'ML-DSA-65 is 4.2x ecrecover — well within Beam\'s 4,500 TPS capacity.' },
            ].map((item, i) => (
              <div key={i} style={{ padding: '8px 0', borderBottom: i < 4 ? '1px solid #ddd' : 'none', display: 'flex', gap: 12 }}>
                <span style={{ ...mono, fontWeight: 700, fontSize: 11, minWidth: 120 }}>{item.label}</span>
                <span style={{ fontSize: 13, color: '#555' }}>{item.desc}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="cell">
          <div style={heading}>Quantum Timeline</div>
          <p style={body}>
            Cryptographically relevant quantum computers are projected in 10-15 years.
            Key migration is slow. The time to start is now.
          </p>
          <div style={{ marginTop: 16, ...mono, fontSize: 11 }}>
            {[
              { year: '1994', event: 'Shor\'s algorithm published — theoretically breaks RSA and ECDSA' },
              { year: '2017', event: 'NIST launches Post-Quantum Cryptography standardization' },
              { year: '2024', event: 'FIPS 204 (ML-DSA) and FIPS 205 (SLH-DSA) finalized' },
              { year: '2026', event: 'PQ_VERIFY precompile built for Beam (this project)' },
              { year: '2027', event: 'Target: Beam mainnet deployment with PQ support' },
              { year: '2030s', event: 'Projected: Cryptographically relevant quantum computers' },
            ].map((item, i) => (
              <div key={i} style={{ padding: '8px 0', borderBottom: i < 5 ? '1px solid #ddd' : 'none', display: 'flex', gap: 12 }}>
                <span style={{ fontWeight: 700, minWidth: 44, color: i === 4 ? 'var(--green)' : '#000' }}>{item.year}</span>
                <span style={{ color: '#555' }}>{item.event}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* ══ How It Works + Performance ════════════════════════════ */}
      <div className="bento-2">
        <div className="cell">
          <div style={heading}>How It Works</div>
          <p style={body}>
            A smart contract sends a <code style={mono}>staticcall</code> to the
            precompile with the public key, signature, and message.
            PQ_VERIFY passes them through liboqs and returns true or false.
          </p>
          <FlowDiagramIcon />
          <p style={{ ...mono, fontSize: 11, color: '#666', marginTop: 12 }}>
            Stateless. No storage writes. No reentrancy risk.
          </p>
        </div>
        <div className="cell cell-clickable" style={{ borderLeft: 'none' }} onClick={() => toggle('performance')}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div style={heading}>Performance</div>
            <span style={toggleIcon(expanded === 'performance')}>+</span>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 0 }}>
            {[
              { val: '105 us', label: 'ML-DSA verify time' },
              { val: '9,504', label: 'Verifications / sec' },
              { val: '133,600', label: 'Gas per ML-DSA verify' },
              { val: '36+', label: 'Tests + fuzz targets' },
            ].map((item, i) => (
              <div key={i} style={{
                padding: '16px', border: '1px solid #000',
                marginTop: i < 2 ? 0 : -1, marginLeft: i % 2 === 0 ? 0 : -1,
              }}>
                <div style={{ fontFamily: 'var(--mono)', fontSize: 22, fontWeight: 700, lineHeight: 1 }}>{item.val}</div>
                <div style={{ fontFamily: 'var(--mono)', fontSize: 10, color: '#666', marginTop: 6, textTransform: 'uppercase', letterSpacing: '0.5px' }}>{item.label}</div>
              </div>
            ))}
          </div>
          <div style={expandHint}>{expanded === 'performance' ? 'Click to collapse' : 'Click for full benchmark log'}</div>
        </div>
      </div>
      {expanded === 'performance' && <PerformanceDetail />}

      {/* ══ Architecture + Integrated + Open Source ═══════════════ */}
      <div className="bento-3">
        <div className="cell">
          <div style={heading}>Architecture</div>
          <LayerStackIcon />
          <div style={{ marginTop: 16, ...mono, fontSize: 11 }}>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}><strong>Layer 3</strong> — Solidity contracts call precompile</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}><strong>Layer 2</strong> — Go precompile with CGo bridge</div>
            <div style={{ padding: '4px 0' }}><strong>Layer 1</strong> — liboqs 0.15 (C, constant-time)</div>
          </div>
        </div>
        <div className="cell" style={{ borderLeft: 'none', borderRight: 'none' }}>
          <div style={heading}>Integrated</div>
          <PuzzleIcon />
          <div style={body}>
            Works with any Solidity contract. Compatible with Subnet-EVM,
            Avalanche tooling, and standard EVM workflows.
          </div>
          <div style={{ marginTop: 16, ...mono, fontSize: 11 }}>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>Subnet-EVM v0.8.0</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>AvalancheGo v1.14.1</div>
            <div style={{ padding: '4px 0', borderBottom: '1px solid #ddd' }}>liboqs 0.15 (Open Quantum Safe)</div>
            <div style={{ padding: '4px 0' }}>EIP-2718 Type 0x50 transactions</div>
          </div>
        </div>
        <div className="cell cell-clickable" onClick={() => toggle('opensource')}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div style={heading}>Open Source</div>
            <span style={toggleIcon(expanded === 'opensource')}>+</span>
          </div>
          <p style={body}>MIT licensed. Built for the Beam Foundation Grant program.</p>
          <div style={{ border: '1px solid #000', padding: '12px 16px', marginTop: 16, ...mono, fontSize: 11 }}>
            <div style={{ fontWeight: 700, marginBottom: 8 }}>Repository</div>
            <span style={{ color: '#000', textDecoration: 'underline' }}>github.com/SAHU-01/pq-beam-verify-precompile</span>
          </div>
          <div style={{ border: '1px solid #000', borderTop: 'none', padding: '12px 16px', ...mono, fontSize: 11 }}>
            <div style={{ fontWeight: 700, marginBottom: 4 }}>Standards</div>
            <div>NIST FIPS 204 / 205</div>
          </div>
          <div style={expandHint}>{expanded === 'opensource' ? 'Click to collapse' : 'Click for repo guide + navigation'}</div>
        </div>
      </div>
      {expanded === 'opensource' && <OpenSourceDetail />}

      {/* ══ Transaction Flow Diagram ═════════════════════════════ */}
      <div style={{ borderBottom: '1px solid #000', padding: '28px' }}>
        <div style={heading}>EIP-2718 Type 0x50 Transaction Flow</div>
        <p style={{ ...body, marginBottom: 20 }}>
          The full path of a post-quantum transaction — from contract construction through liboqs verification.
        </p>
        <TransactionFlowDiagram />
      </div>

      {/* ══ On-Chain Proof ════════════════════════════════════════ */}
      <div style={{ borderBottom: '1px solid #000' }}>
        <div className="cell" style={{ borderBottom: 'none' }}>
          <div style={heading}>Proven On-Chain</div>
          <p style={{ ...body, marginBottom: 0 }}>Deployed and tested on a local Beam subnet. Real signatures — not simulated.</p>
        </div>
        <div className="bento-3">
          <div className="cell" style={{ borderTop: 'none', borderLeft: '3px solid var(--green)' }}>
            <div style={{ ...mono, fontSize: 10, letterSpacing: '2px', textTransform: 'uppercase', color: 'var(--green)', fontWeight: 700, marginBottom: 8 }}>VERIFIED ON-CHAIN</div>
            <div style={{ ...mono, fontSize: 11, wordBreak: 'break-all', marginBottom: 12, color: '#666' }}>TX: 0x5bc2aff8...52febc</div>
            <div style={{ ...mono, fontSize: 11, lineHeight: 1.8 }}>
              <div>Message: "Hello Beam! Post-quantum signatures are live on-chain."</div>
              <div>Pubkey: 1,952 bytes (ML-DSA-65)</div>
              <div>Signature: 3,309 bytes</div>
              <div style={{ color: 'var(--green)', fontWeight: 700, marginTop: 4 }}>Result: valid = true</div>
            </div>
          </div>
          <div className="cell" style={{ borderTop: 'none', borderLeft: '3px solid var(--red)' }}>
            <div style={{ ...mono, fontSize: 10, letterSpacing: '2px', textTransform: 'uppercase', color: 'var(--red)', fontWeight: 700, marginBottom: 8 }}>REJECTED ON-CHAIN</div>
            <div style={{ ...mono, fontSize: 11, wordBreak: 'break-all', marginBottom: 12, color: '#666' }}>TX: 0x9cb8d4f6...10907</div>
            <div style={{ ...mono, fontSize: 11, lineHeight: 1.8 }}>
              <div>Same pubkey and message.</div>
              <div>1 byte flipped in signature.</div>
              <div style={{ color: 'var(--red)', fontWeight: 700, marginTop: 4 }}>Result: valid = false</div>
            </div>
          </div>
          <div className="cell" style={{ borderTop: 'none' }}>
            <div style={{ ...mono, fontSize: 10, letterSpacing: '2px', textTransform: 'uppercase', color: '#666', fontWeight: 700, marginBottom: 8 }}>CHAIN DETAILS</div>
            <div style={{ ...mono, fontSize: 11, lineHeight: 2 }}>
              <div>Chain ID: 13337</div>
              <div>VM: Subnet-EVM + PQ_VERIFY</div>
              <div>Block: #2</div>
              <div>Gas used: 250,146</div>
              <div>Precompile gas: 140,072</div>
            </div>
          </div>
        </div>
      </div>

      {/* ══ Algorithm Comparison ══════════════════════════════════ */}
      <div style={{ borderBottom: '1px solid #000', padding: '28px' }}>
        <div style={heading}>Algorithm Comparison</div>
        <div style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontFamily: 'var(--mono)', fontSize: 12 }}>
            <thead><tr>
              {['Algorithm', 'Standard', 'Public Key', 'Signature', 'Verify', 'Gas', 'Security'].map(h => (
                <th key={h} style={{ textAlign: 'left', padding: '10px 12px', borderBottom: '2px solid #000', fontSize: 10, textTransform: 'uppercase', letterSpacing: '1px', fontWeight: 700 }}>{h}</th>
              ))}
            </tr></thead>
            <tbody>
              {[
                { alg: 'ML-DSA-65', std: 'FIPS 204', pk: '1,952 B', sig: '3,309 B', time: '105 us', gas: '133,600', sec: 'Level 3' },
                { alg: 'SLH-DSA-128s', std: 'FIPS 205', pk: '32 B', sig: '7,856 B', time: '432 us', gas: '523,600', sec: 'Level 1' },
                { alg: 'ECDSA (current)', std: '--', pk: '33 B', sig: '65 B', time: '~25 us', gas: '3,000', sec: null },
              ].map((r, i) => (
                <tr key={i} style={i === 2 ? { background: '#f8f8f8' } : {}}>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd', fontWeight: 700 }}>{r.alg}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>{r.std}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>{r.pk}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>{r.sig}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>{r.time}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>{r.gas}</td>
                  <td style={{ padding: '10px 12px', borderBottom: '1px solid #ddd' }}>
                    {r.sec === null
                      ? <span style={{ ...mono, fontSize: 10, padding: '2px 8px', border: '1px solid var(--red)', color: 'var(--red)', fontWeight: 700 }}>BROKEN BY SHOR'S</span>
                      : r.sec}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* ══ Roadmap — CLICKABLE PHASES ════════════════════════════ */}
      <div className="bento-3">
        {[
          { phase: 'Phase 1', title: 'Precompile', status: 'IN PROGRESS', key: 'phase1',
            items: ['PQ_VERIFY precompile (MVP live)', 'Beam testnet fork deployment', 'Gas calibration on validators', 'ERC-4337 + key rotation contracts'], done: false },
          { phase: 'Phase 2', title: 'SDK + Native Support', status: 'NEXT', key: 'phase2',
            items: ['TypeScript SDK', 'Key management', 'Type 0x50 in validators', 'Auto PQ key creation'], done: false },
          { phase: 'Phase 3', title: 'Audit + Mainnet', status: 'PLANNED', key: 'phase3',
            items: ['Security audit (CGo)', 'Validator benchmarks', 'Mainnet deployment', 'Migration toolkit'], done: false },
        ].map((p, i) => (
          <div key={i} className="cell cell-clickable" onClick={() => toggle(p.key)}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
              <span style={{ ...mono, fontSize: 10, letterSpacing: '2px', textTransform: 'uppercase', color: '#666' }}>{p.phase}</span>
              <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                <span style={{ ...mono, fontSize: 9, padding: '2px 8px', border: `1px solid ${p.status === 'IN PROGRESS' ? '#d97706' : '#000'}`, color: p.status === 'IN PROGRESS' ? '#d97706' : '#000', fontWeight: 700, letterSpacing: '1px' }}>{p.status}</span>
                <span style={toggleIcon(expanded === p.key)}>+</span>
              </div>
            </div>
            <div style={{ ...mono, fontSize: 13, fontWeight: 700, marginBottom: 12 }}>{p.title}</div>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {p.items.map((item, j) => (
                <li key={j} style={{ ...mono, fontSize: 11, padding: '4px 0', color: '#333', display: 'flex', gap: 8 }}>
                  <span style={{ color: p.done ? 'var(--green)' : '#999' }}>{p.done ? '+' : '-'}</span>{item}
                </li>
              ))}
            </ul>
            <div style={expandHint}>{expanded === p.key ? 'Click to collapse' : 'Click for milestones + details'}</div>
          </div>
        ))}
      </div>
      {expanded === 'phase1' && <PhaseDetail phase={1} />}
      {expanded === 'phase2' && <PhaseDetail phase={2} />}
      {expanded === 'phase3' && <PhaseDetail phase={3} />}

      {/* ══ PQ Transaction Type ═══════════════════════════════════ */}
      <div className="bento-2">
        <div className="terminal-cell">
          <div style={{ color: '#888', marginBottom: 8, fontSize: 11 }}>// EIP-2718 Typed Transaction Envelope</div>
          <div style={{ color: 'var(--terminal-green)', marginBottom: 12 }}>Type 0x50</div>
          <div style={{ color: '#ccc', lineHeight: 2 }}>
            0x50 || RLP<span style={{ color: '#888' }}>([</span><br />
            {'  '}chainId,<br />
            {'  '}nonce,<br />
            {'  '}gasPrice, gasLimit,<br />
            {'  '}to, value, data,<br />
            {'  '}<span style={{ color: 'var(--terminal-green)' }}>pqAlgorithm</span>,{'   '}<span style={{ color: '#888' }}>// uint8</span><br />
            {'  '}<span style={{ color: 'var(--terminal-green)' }}>pqPublicKey</span>,{'   '}<span style={{ color: '#888' }}>// bytes</span><br />
            {'  '}<span style={{ color: 'var(--terminal-green)' }}>pqSignature</span>{'   '}<span style={{ color: '#888' }}>// bytes</span><br />
            <span style={{ color: '#888' }}>])</span>
          </div>
        </div>
        <div className="cell">
          <div style={heading}>PQ Transaction Type</div>
          <div style={{ ...mono, fontSize: 11, lineHeight: 2 }}>
            <div style={{ padding: '6px 0', borderBottom: '1px solid #ddd' }}><strong>Type byte:</strong> 0x50 (ASCII 'P')</div>
            <div style={{ padding: '6px 0', borderBottom: '1px solid #ddd' }}><strong>Address:</strong> keccak256(pqPubKey)[12:32]</div>
            <div style={{ padding: '6px 0', borderBottom: '1px solid #ddd' }}><strong>Replay protection:</strong> Chain ID + type prefix</div>
            <div style={{ padding: '6px 0' }}><strong>Backward compatible:</strong> ECDSA unchanged</div>
          </div>
        </div>
      </div>

      {/* ══ Test Coverage ════════════════════════════════════════ */}
      <div style={{ borderBottom: '1px solid #000', padding: '28px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
          <div style={heading}>Test Coverage</div>
          <button onClick={() => setShowRun(!showRun)} style={{
            fontFamily: 'var(--mono)', fontSize: 11, padding: '6px 14px', border: '1px solid #000',
            background: '#fff', cursor: 'pointer', fontWeight: 700, letterSpacing: '1px', textTransform: 'uppercase',
          }}>{showRun ? 'Hide' : 'Run'} Commands</button>
        </div>
        <div className="bento-stats">
          {[
            { suite: 'pkg/pqcrypto', tests: 8, desc: 'Keygen, sign, verify, tampering' },
            { suite: 'pkg/pqverify', tests: 12, desc: 'ABI, gas, valid/invalid + 2 fuzz targets' },
            { suite: 'test/e2e', tests: 7, desc: 'Full flow, dual alg, 1MB message' },
            { suite: 'cmd/benchmark', tests: 4, desc: 'Go benchmarks, both algorithms' },
            { suite: 'Solidity', tests: 6, desc: 'PQAccount (ERC-4337), PQKeyRotation' },
          ].map((t, i) => (
            <div key={i} style={{ border: '1px solid #000', padding: '16px', marginLeft: i > 0 ? -1 : 0 }}>
              <div style={{ ...mono, fontWeight: 700, fontSize: 12 }}>{t.suite}</div>
              <div style={{ fontFamily: 'var(--mono)', fontSize: 24, fontWeight: 700, margin: '8px 0' }}>{t.tests}</div>
              <div style={{ ...mono, fontSize: 10, color: '#666' }}>{t.desc}</div>
            </div>
          ))}
        </div>
        {showRun && (
          <pre style={{ background: 'var(--terminal-bg)', color: 'var(--terminal-text)', fontFamily: 'var(--mono)', fontSize: 11, padding: '16px 20px', lineHeight: 1.7, border: '1px solid #000', overflow: 'auto', marginTop: 16 }}>
{`# Unit tests
CGO_ENABLED=1 go test ./pkg/... -v

# End-to-end tests
CGO_ENABLED=1 go test ./test/ -v

# Benchmarks
CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=. -benchmem

# On-chain demo (local subnet must be running)
RPC_URL=<rpc-url> ./scripts/demo_onchain.sh`}
          </pre>
        )}
      </div>

      {/* ══ Security ═════════════════════════════════════════════ */}
      <div className="bento-3">
        {[
          { title: 'Stateless', desc: 'No storage reads or writes. Pure function. No reentrancy, no state corruption.' },
          { title: 'Gas-Metered', desc: 'Cost proportional to computation with 10x safety margin. Prevents DoS.' },
          { title: 'NIST Standard', desc: 'Both algorithms are finalized FIPS standards (2024). Production-grade.' },
        ].map((item, i) => (
          <div key={i} className="cell"><div style={heading}>{item.title}</div><p style={body}>{item.desc}</p></div>
        ))}
      </div>
      <div className="bento-3" style={{ borderTop: 'none' }}>
        {[
          { title: 'Dual Algorithm', desc: 'Hash-based fallback if lattice assumptions are ever weakened.' },
          { title: 'Constant-Time', desc: 'liboqs designed for constant-time. CGo boundary under Phase 3 audit.' },
          { title: 'Ephemeral Tests', desc: 'All test data generated fresh per run. No hardcoded keys.' },
        ].map((item, i) => (
          <div key={i} className="cell" style={{ borderTop: 'none' }}>
            <div style={heading}>{item.title}</div><p style={body}>{item.desc}</p>
          </div>
        ))}
      </div>

      {/* ══ References ═══════════════════════════════════════════ */}
      <div style={{ padding: '28px', borderBottom: '1px solid #000' }}>
        <div style={heading}>References</div>
        <div className="ref-columns">
          {[
            { id: '1', title: 'FIPS 204 — ML-DSA', org: 'NIST, 2024', url: 'https://csrc.nist.gov/pubs/fips/204/final' },
            { id: '2', title: 'FIPS 205 — SLH-DSA', org: 'NIST, 2024', url: 'https://csrc.nist.gov/pubs/fips/205/final' },
            { id: '3', title: 'CRYSTALS-Dilithium', org: 'Ducas et al., 2021', url: 'https://pq-crystals.org/dilithium/' },
            { id: '4', title: 'SPHINCS+', org: 'Aumasson et al., 2022', url: 'https://sphincs.org/' },
            { id: '5', title: 'Quantum Computing & Finance', org: 'IMF, 2024', url: 'https://www.imf.org/en/Publications/fintech-notes/Issues/2024/03/20/Quantum-Computing-and-the-Financial-System-Spooky-Action-at-a-Distance-546082' },
            { id: '6', title: "Shor's Algorithm", org: 'Peter W. Shor, 1994', url: 'https://arxiv.org/abs/quant-ph/9508027' },
            { id: '7', title: 'Open Quantum Safe — liboqs', org: 'OQS Project, 2024', url: 'https://github.com/open-quantum-safe/liboqs' },
            { id: '8', title: 'EIP-2718: Typed Transactions', org: 'Ethereum Foundation, 2020', url: 'https://eips.ethereum.org/EIPS/eip-2718' },
            { id: '9', title: 'Subnet-EVM Precompiles', org: 'Ava Labs, 2024', url: 'https://docs.avax.network/build/subnet/upgrade/customize-a-subnet#precompiles' },
          ].map((ref) => (
            <div key={ref.id} style={{ display: 'flex', gap: 8, padding: '6px 0', breakInside: 'avoid' as const, alignItems: 'baseline' }}>
              <span style={{ ...mono, fontSize: 10, color: '#999', minWidth: 20 }}>[{ref.id}]</span>
              <div>
                <a href={ref.url} target="_blank" rel="noopener noreferrer" style={{ ...mono, fontSize: 11, color: '#000', textDecoration: 'underline' }}>{ref.title}</a>
                <div style={{ ...mono, fontSize: 10, color: '#999' }}>{ref.org}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* ══ Footer ═══════════════════════════════════════════════ */}
      <footer style={{ padding: '32px 0', textAlign: 'center', fontFamily: 'var(--mono)', fontSize: 11, color: '#999' }}>
        <div>
          Built for the{' '}
          <a href="https://www.onbeam.com/" target="_blank" rel="noopener noreferrer" style={{ color: '#000', textDecoration: 'underline' }}>Beam Foundation Grant</a>
          {' — '}Post-Quantum Signing Infrastructure
        </div>
        <div style={{ marginTop: 8 }}>
          <a href="https://github.com/SAHU-01/pq-beam-verify-precompile" target="_blank" rel="noopener noreferrer" style={{ color: '#000', textDecoration: 'underline' }}>GitHub</a>
          {' / '}
          <a href="https://csrc.nist.gov/pubs/fips/204/final" target="_blank" rel="noopener noreferrer" style={{ color: '#000', textDecoration: 'underline' }}>FIPS 204</a>
          {' / '}
          <a href="https://csrc.nist.gov/pubs/fips/205/final" target="_blank" rel="noopener noreferrer" style={{ color: '#000', textDecoration: 'underline' }}>FIPS 205</a>
          {' / MIT License'}
        </div>
      </footer>
    </div>
  )
}

export default App
