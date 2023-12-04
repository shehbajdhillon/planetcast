import {
  ArrowUpFromDot,
  CircleDollarSign,
  DollarSign,
  Globe2Icon,
  GlobeIcon,
  LanguagesIcon,
  TimerReset,
  TrendingDownIcon,
  Volume2Icon,
  ImportIcon,
} from 'lucide-react';

import {
  Box,
  HStack,
  Heading,
  Stack,
  Text,
  useBreakpointValue,
} from "@chakra-ui/react";

const BenefitsSection: React.FC = () => {
  const iconSize = useBreakpointValue({ base: '40px', md: '60px' })
  const iconSizeBig = useBreakpointValue({ base: "75px", sm: "120px", md: "140px" });

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >

      <Heading
        size={{ base: '2xl', md: '3xl' }}
        fontWeight={'medium'}
        textAlign={"left"}
        w={"full"}
      >
        <HStack>
          <Text>10x your {' '}
            <Text
              as={"span"}
              bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
              bgClip='text'
            >
              reach
            </Text>
          </Text>
          <ArrowUpFromDot size={iconSize} />
          <GlobeIcon size={iconSize} />
        </HStack>
        <HStack mt={{ md: "8px" }}>
          <Text>1/10th the {' '}
            <Text
              as={"span"}
              bgGradient={'linear(to-tr, #01CF00, #90DD00)'}
              bgClip='text'
            >
              cost
            </Text>
          </Text>
          <TrendingDownIcon size={iconSize} />
          <DollarSign size={iconSize} />
        </HStack>
      </Heading>

      <Box
        mb={{ base: "auto", md: "0px" }}
        w="full"
        alignItems={"left"}
        justifyContent={"left"}
        display={"flex"}
        flexDir={"column"}
      >
        <Stack direction={{ base: "column", md: "row" }} mt={{ base: "25px" }} spacing={{ base: "25px", md: "0px" }}>

          <Stack w="full" direction={{ base: "row", md: "column" }}>
            <Heading
              fontWeight={'normal'}
              textAlign={"left"}
              w={"full"}
              size={{ base: "sm", sm: "lg" }}
            >
              {"Engage listeners everywhere around the globe. Over 80% of the world's population does not speak English."}
            </Heading>
            <HStack w="full">
              <Globe2Icon size={iconSizeBig} strokeWidth={"0.75px"}/>
              <LanguagesIcon size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

          <Stack w="full" direction={{ base: "row-reverse", md: "column" }}>
            <Heading
              textAlign={"left"}
              w={"full"}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
            >
              Professional dubbing takes longer and breaks your bank.
              Our AI powered tools help you save time and money.
            </Heading>
            <HStack w="full">
              <CircleDollarSign size={iconSizeBig} strokeWidth={"0.75px"} />
              <TimerReset size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

          <Stack w="full" direction={{ base: "row", md: "column" }}>
            <Heading
              textAlign={"left"}
              w={"full"}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
            >
              Preserve original voices in every translation.
              Our AI powered tools preserve the original voice of the speaker in all the dubbings.
            </Heading>
            <HStack w="full">
              <Volume2Icon size={iconSizeBig} strokeWidth={"0.75px"} />
              <ImportIcon size={iconSizeBig} strokeWidth={"0.75px"} />
            </HStack>
          </Stack>

        </Stack>
      </Box>
    </Stack>
  );
};

export default BenefitsSection;
