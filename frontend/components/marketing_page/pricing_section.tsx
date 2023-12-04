import { HStack, Heading, Stack, Text, useColorModeValue } from "@chakra-ui/react";
import { useState } from "react";
import Button from "../button";
import PricingComponent from "./pricing_component";

const PricingSection = () => {

  const [annualPricing, setAnnualPricing] = useState(true);

  const bgColor = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={"left"}
        w={"full"}
      >
        Start dubbing {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          today
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={"left"}
        w={"full"}
      >
        Start for Free. No Credit Card Required.
      </Heading>

      <HStack borderWidth={"1px"} p="5px" my="40px" rounded={"md"}>
        <Button
          onClick={() => setAnnualPricing(false)}
          bgColor={!annualPricing ? bgColor : textColor}
          textColor={!annualPricing ? textColor : bgColor}
          bgGradient={!annualPricing ? 'linear(to-tl, #007CF0, #01DFD8)' : ''}
          _hover={{
            backgroundColor: !annualPricing ? bgColor : textColor,
            textColor: annualPricing ? bgColor : textColor,
          }}
        >
          Monthly
        </Button>
        <Button
          onClick={() => setAnnualPricing(true)}
          bgColor={annualPricing ? bgColor : textColor}
          textColor={annualPricing ? textColor : bgColor}
          bgGradient={annualPricing ? 'linear(to-tl, #007CF0, #01DFD8)' : ''}
          _hover={{
            backgroundColor: annualPricing ? bgColor : textColor,
            textColor: !annualPricing ? bgColor : textColor,
          }}
        >
          Annual
        </Button>
      </HStack>

      <PricingComponent annualPricing={annualPricing} marketingPage={true} />
    </Stack>
  );
};

export default PricingSection;
