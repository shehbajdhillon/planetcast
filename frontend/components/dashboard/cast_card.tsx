import {
  Box,
  HStack,
  Spacer,
  Text,
  Button,
  useColorModeValue,
} from "@chakra-ui/react";

interface CastCardProps {
  title: string;
  status: "DRAFT" | "PROCESSING" | "DONE";
};

const CastCard: React.FC<CastCardProps> = (props) => {

  const { title, status } = props;

  const statusColorScheme = status === "DRAFT" ? "gray" : status === "PROCESSING" ? "blue" : "green"

  return (
    <Box
      borderWidth={"1px"}
      p={6}
      maxW="400px"
      minW="270px"
      rounded={"lg"}
      _hover={{
        borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
        boxShadow: 'lg',
        bg: useColorModeValue('white', 'whiteAlpha.100'),
      }}
      cursor={"pointer"}
    >
      <HStack>
        <Text
          textTransform="capitalize"
          maxW="400px"
          fontWeight={700}
          fontSize={'lg'}
          letterSpacing={1.1}
          noOfLines={1}
        >
          {title}
        </Text>
      </HStack>
      <Spacer />
      <HStack pt="20px">
        <Button
          borderWidth="1px"
          size={'xs'}
          textTransform="capitalize"
          fontWeight="medium"
          alignContent="right"
          colorScheme={statusColorScheme}
          pointerEvents={"none"}
        >
          {status}
        </Button>
      </HStack>
    </Box>
  );
};

export default CastCard;
