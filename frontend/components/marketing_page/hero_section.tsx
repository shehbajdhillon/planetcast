import {
  Box,
  Button,
  HStack,
  Heading,
  Stack,
  Text,
  useBreakpointValue,
  useColorModeValue
} from "@chakra-ui/react";
import Link from "next/link";
import { DarkModeGradientLogo, LightModeGradientLogo } from "../logo";

const HeroSection: React.FC = () => {

  const whiteBlack = useColorModeValue("white", "black");
  const blackWhite = useColorModeValue("black", "white");

  const logoSize = useBreakpointValue({ base: 200, md: 250, lg: 400 });
  const darkModeOn = useColorModeValue(false, true);

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      direction={{ base: "column-reverse", md: "row" }}
      maxW={"1400px"}
      w="full"
    >
      <Box
        mb={{ base: "auto", md: "0px" }}
        w="full"
        maxW={{ md: "75%" }}
        alignItems={{ base: "center", md: "left" }}
        justifyContent={{ base: "center", md: "left" }}
        display={"flex"}
        flexDir={"column"}
      >
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Dub
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Translate
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Broadcast
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Content Across the {' '}
          <Text
            as={"span"}
            bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
            bgClip='text'
          >
            Planet
          </Text>
        </Heading>
        <HStack w={{ md: "full" }} pt="10px">
          <Link href={'/dashboard'}>
            <Button
              size={"lg"}
              backgroundColor={blackWhite}
              textColor={whiteBlack}
              borderColor={whiteBlack}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: whiteBlack,
                textColor: whiteBlack,
                bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
              }}
            >
              Start for Free
            </Button>
          </Link>
          <Link href={"#usecases"}>
            <Button
              size={"lg"}
              variant={"outline"}
              _hover={{
                backgroundColor: blackWhite,
                textColor: whiteBlack,
                bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
              }}
            >
              Read More
            </Button>
          </Link>
        </HStack>
        <Text w={{ md: "full" }} pt="2px" pl="1px">
          No Credit Card Required
        </Text>
      </Box>
      <Box maxW={{ base: "200px", md: "25%" }} mt={{ base:"auto", md: "0px" }}>
        {
          darkModeOn ?
          <DarkModeGradientLogo height={logoSize || 0} width={logoSize || 0} />
          :
          <LightModeGradientLogo height={logoSize || 0} width={logoSize || 0} />
        }
      </Box>
    </Stack>
  );
};

export default HeroSection;
