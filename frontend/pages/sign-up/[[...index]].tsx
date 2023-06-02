import { Box } from "@chakra-ui/react";
import { SignUp } from "@clerk/nextjs";

export default function Page() {
  return (
    <Box py={{ base: '150px', md: '200px' }} alignItems={"center"} display={"flex"} justifyContent={"center"}>
      <SignUp />
    </Box>
  );
}
