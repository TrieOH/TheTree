import { motion, AnimatePresence } from "motion/react";
import { cn } from "@/shared/lib/utils";

interface PropsI {
  message?: string;
  className?: string;
}

const FormError = ({ message, className }: PropsI) => {
  return (
    <AnimatePresence>
      {message && (
        <motion.p
          initial={{ opacity: 0, height: 0, marginTop: 0 }}
          animate={{ opacity: 1, height: "auto", marginTop: 6 }}
          exit={{ opacity: 0, height: 0, marginTop: 0 }}
          transition={{ duration: 0.2, ease: "easeOut" }}
          className={cn(
            "text-[11px] font-medium text-destructive tracking-tight leading-none px-1 overflow-hidden",
            className
          )}
        >
          {message}
        </motion.p>
      )}
    </AnimatePresence>
  );
};

export default FormError;
