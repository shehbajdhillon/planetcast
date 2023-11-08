import { Box, Text, Button, VStack, HStack, Icon, useColorModeValue, Badge, Stack, Circle, Spacer, Heading } from "@chakra-ui/react";
import { CheckIcon } from "lucide-react";
import Link from "next/link";

const PricingComponent = () => {
  const priceColor= useColorModeValue("zinc.600", "zinc.400");
  const cardBgColor = useColorModeValue("white", "black");

  const priceMap: Record<string, any> = {
    'Starter': {
      'price': 57,
      'features': [
        '30 minutes of dubbing included',
        '$1.97 per additional minute of dubbing',
        'Dub videos to 28+ Languagues',
      ],
    },
    'Pro': {
      'price': 137,
      'features': [
        '100 minutes of dubbing included',
        '$1.47 per additional minute of dubbing',
        'Dub videos to 28+ Languagues',
      ],
    },
    'Business': {
      'price': 497,
      'features': [
        '500 minutes of dubbing included',
        '$0.97 per additional minute of dubbing',
        'Dub videos to 28+ Languagues',
        'Descript, Rumble, YouTube Integration',
        'API Access',
        "CEO's Phone number",
      ],
    }
  };

  const bgColor = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

  return (
    <Stack direction={{ base: "column", md: "row" }} w="full" spacing={"80px"} px="16px">
      {/* Repeating pattern for Cards, replace as needed */}
      {Object.keys(priceMap).map((tier)  => (
        <Box
          flex="1"
          p={6}
          bg={cardBgColor}
          shadow="lg"
          borderRadius="md"
          borderWidth={1}
          position="relative"
          key={tier}
        >
          {tier === "Pro" && (
            <Badge
              position="absolute"
              top="-10px"
              left="50%"
              transform="translateX(-50%)"
              backgroundColor={bgColor}
              textColor={textColor}
              bgGradient='linear(to-tl, #007CF0, #01DFD8)'
            >
              Popular
            </Badge>
          )}
          <VStack spacing={4} align="center" h="full">

            <Heading fontWeight={"semibold"}>{tier}</Heading>
            <Text fontSize="4xl" fontWeight="normal" color={priceColor}>
              ${priceMap[tier].price}/month
            </Text>
            <VStack align="start" spacing={2}>
              {priceMap[tier].features.map((detail: string, idx: number) => (
                <HStack key={idx}>
                  <Circle bg="green.500" borderRadius="full" p={1}>
                    <Icon as={CheckIcon} color="white" />
                  </Circle>
                  <Text>{detail}</Text>
                </HStack>
              ))}
            </VStack>
            <Spacer />

            <Link href={'/dashboard'}>
              <Button
                size={"lg"}
                backgroundColor={bgColor}
                textColor={textColor}
                borderColor={textColor}
                borderWidth={"1px"}
                bgGradient={tier === "Pro" ? 'linear(to-tl, #007CF0, #01DFD8)' : ''}
                _hover={{
                  backgroundColor: tier === "Pro" ? bgColor : textColor,
                  textColor: textColor,
                  bgGradient: tier === "Pro" ? '' : 'linear(to-tl, #007CF0, #01DFD8)'
                }}
              >
                Start for Free
              </Button>
            </Link>
          </VStack>
        </Box>
      ))}
    </Stack>
  )
}

export default PricingComponent;

