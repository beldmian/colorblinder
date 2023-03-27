import React from 'react'

type FilterProps = {
    color: string,
    color_opacity: number,
    contrast: number,
    saturation: number,
}

export default function Filter({ color, color_opacity, contrast, saturation }: FilterProps) {
  return (
    <svg width="0" height="0">
        <defs>
            <filter id="video_filter" x="0" y="0" width="100%" height="100%" filterUnits="objectBoundingBox" primitiveUnits="userSpaceOnUse" colorInterpolationFilters="linearRGB">
                <feFlood floodColor={color} floodOpacity={color_opacity} result="flood"/>
                <feComposite in="flood" in2="SourceAlpha" operator="in" result="composite"/>
                <feBlend mode="luminosity" in="SourceGraphic" in2="composite" result="blend3"/>
                <feComponentTransfer in="blend3" result="componentTransfer2">
                    <feFuncR type="linear" slope={contrast} intercept="0"/>
                    <feFuncG type="linear" slope={contrast} intercept="0"/>
                    <feFuncB type="linear" slope={contrast} intercept="0"/>
                    <feFuncA type="identity"/>
                </feComponentTransfer>
                <feColorMatrix type="saturate" values={saturation.toString()} in="componentTransfer2" result="colormatrix2"/>
            </filter>
        </defs>
    </svg>
  )
}
