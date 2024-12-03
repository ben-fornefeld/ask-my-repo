import { clsx, type ClassValue } from "clsx";
import { Variants } from "framer-motion";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const brutalistMotion: Variants = {
  hidden: {
    opacity: 0,
    scale: 0.95,
    y: 10,
    rotate: -2,
    transition: {
      type: "spring",
      damping: 20,
      stiffness: 300,
    },
  },
  visible: {
    opacity: 1,
    scale: 1,
    y: 0,
    rotate: 0,
    transition: {
      type: "spring",
      damping: 15,
      stiffness: 300,
      mass: 0.8,
    },
  },
  hover: {
    scale: 1.02,
    rotate: 1,
    transition: {
      type: "spring",
      damping: 10,
      stiffness: 300,
    },
  },
  tap: {
    scale: 0.98,
    rotate: -1,
    transition: {
      type: "spring",
      damping: 10,
      stiffness: 300,
    },
  },
};

export const brutalistSlideMotion: Variants = {
  hidden: {
    opacity: 0,
    x: -40,
    rotate: -1,
    scale: 0.98,
    transition: {
      type: "spring",
      damping: 20,
      stiffness: 500,
    },
  },
  visible: {
    opacity: 1,
    x: 0,
    rotate: 0,
    scale: 1,
    transition: {
      type: "spring",
      damping: 25,
      stiffness: 500,
      mass: 0.8,
      delay: 0.05,
    },
  },
  exit: {
    opacity: 0,
    x: 40,
    rotate: 1,
    scale: 0.98,
    transition: {
      type: "spring",
      damping: 20,
      stiffness: 500,
    },
  },
};
