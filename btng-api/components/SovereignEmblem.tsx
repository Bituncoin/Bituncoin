'use client'

import { useState, useEffect } from 'react'

interface SovereignEmblemProps {
  size?: number
  className?: string
  trigger?: boolean
  onTriggerComplete?: () => void
}

export default function SovereignEmblem({ 
  size = 40, 
  className = '',
  trigger = false,
  onTriggerComplete
}: SovereignEmblemProps) {
  const [isTriggered, setIsTriggered] = useState(false)

  useEffect(() => {
    if (trigger) {
      setIsTriggered(true)
      const timer = setTimeout(() => {
        setIsTriggered(false)
        onTriggerComplete?.()
      }, 600) // Match emblem-trigger-burst animation duration
      
      return () => clearTimeout(timer)
    }
  }, [trigger, onTriggerComplete])

  const emblemClasses = [
    'sovereign-emblem',
    isTriggered && 'emblem-trigger',
    className
  ].filter(Boolean).join(' ')

  return (
    <div className={emblemClasses} style={{ width: size, height: size }}>
      <svg 
        width="100%" 
        height="100%" 
        viewBox="0 0 200 200" 
        xmlns="http://www.w3.org/2000/svg"
      >
        {/* Golden Shield */}
        <defs>
          <linearGradient id="shield-gradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#F4E4B7" />
            <stop offset="50%" stopColor="#D4AF37" />
            <stop offset="100%" stopColor="#B8941F" />
          </linearGradient>
          <filter id="emblem-inner-glow">
            <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
            <feMerge>
              <feMergeNode in="coloredBlur"/>
              <feMergeNode in="SourceGraphic"/>
            </feMerge>
          </filter>
        </defs>
        
        <path 
          d="M100 10 L160 40 L160 120 C160 160 100 190 100 190 C100 190 40 160 40 120 L40 40 Z" 
          fill="url(#shield-gradient)" 
          stroke="#B8941F" 
          strokeWidth="3"
          filter="url(#emblem-inner-glow)"
        />
        
        {/* Trust Symbol (Interlocking Rings) */}
        <circle 
          cx="85" 
          cy="80" 
          r="25" 
          fill="none" 
          stroke="#1A4B8C" 
          strokeWidth="4"
          opacity="0.9"
        />
        <circle 
          cx="115" 
          cy="80" 
          r="25" 
          fill="none" 
          stroke="#1A4B8C" 
          strokeWidth="4"
          opacity="0.9"
        />
        
        {/* BTNG Text */}
        <text 
          x="100" 
          y="140" 
          fontFamily="Arial, sans-serif" 
          fontSize="24" 
          fontWeight="bold" 
          textAnchor="middle" 
          fill="#0A0A0A"
        >
          BTNG
        </text>
        
        {/* Tagline */}
        <text 
          x="100" 
          y="160" 
          fontFamily="Arial, sans-serif" 
          fontSize="10" 
          textAnchor="middle" 
          fill="#1A4B8C"
          opacity="0.8"
        >
          BUILDING TRUST
        </text>
      </svg>
    </div>
  )
}
