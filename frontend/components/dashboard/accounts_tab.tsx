import useWindowDimensions from "@/hooks/useWindowDimensions";
import {
  Box,
  Heading,
  Divider,
  useColorModeValue,
  Stack,
  Text,
  HStack,
  Badge,
  Spacer,
  useDisclosure,
  useBreakpointValue,
  Drawer,
  DrawerOverlay,
  DrawerContent,
  DrawerCloseButton,
  DrawerBody,
} from "@chakra-ui/react";
import { LandmarkIcon } from "lucide-react";
import { useEffect, useState } from "react";

const AccountCardData = [
  {
    id: 42,
    accountName: "Wells Fargo Everyday Checking",
    hoursUpdatedAgo: 9,
    amount: 3455,
    percentageChange: 103,
    isProfit: true,
  },
  {
    id: 56,
    accountName: "Chase Checking",
    hoursUpdatedAgo: 3,
    amount: 4523,
    percentageChange: 43,
    isProfit: false,
  },
  {
    id: 52,
    accountName: "Bank of America",
    hoursUpdatedAgo: 3,
    amount: 4523,
    percentageChange: 43,
    isProfit: false,
  },
  {
    id: 1230,
    accountName: "Bank Of India",
    hoursUpdatedAgo: 3,
    amount: 4523,
    percentageChange: 43,
    isProfit: false,
  },
  {
    id: 3130,
    accountName: "Credit Suisse",
    hoursUpdatedAgo: 3,
    amount: 4523,
    percentageChange: 43,
    isProfit: false,
  },
];


interface AccountCardProps {
  accountName: string;
  hoursUpdatedAgo: number;
  amount: number;
  percentageChange: number;
  isProfit: boolean;
  onClick: () => void;
};


interface AccountDrawerProps {
  accountId: number;
};

const AccountDrawer: React.FC<AccountDrawerProps> = ({ accountId }) => {
  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");
  return (
    <Stack
      display={"flex"}
      w="full"
      h="full"
      overflow={"auto"}
      spacing={"20px"}
    >
      <Box>
        <Box
          h="250px"
          w="full"
          borderColor={dividerColor}
          borderWidth={"1px"}
          borderRadius={"md"}
        />
        <Stack spacing={"10px"}>
          <Heading>
            {AccountCardData.find(card => card.id === accountId)?.accountName}
          </Heading>
          {AccountCardData.filter((card) => card.id === accountId).map((card, idx) => {
            return (
              <AccountCard
                key={idx}
                accountName={card.accountName}
                hoursUpdatedAgo={card.hoursUpdatedAgo}
                amount={card.amount}
                percentageChange={card.percentageChange}
                isProfit={card.isProfit}
                onClick={() => {}}
              />
            );
          })}
        </Stack>
      </Box>
    </Stack>
  );
};

const AccountCard: React.FC<AccountCardProps> = (props) => {

  const {
    accountName,
    hoursUpdatedAgo,
    amount,
    percentageChange,
    isProfit,
    onClick,
  } = props;

  return (
    <HStack
      onClick={onClick}
      borderWidth={"1px"}
      borderRadius={"lg"}
      px="10px"
      justifyContent={"space-between"}
      _hover={{
        backgroundColor: useColorModeValue("gray.100", "whiteAlpha.200"),
        cursor: "pointer",
      }}
    >
      <Stack
        spacing={"0px"}
        maxW={"30%"}
        w="full"
      >
        <Text
          whiteSpace={"nowrap"}
          textOverflow={"ellipsis"}
          overflow={'hidden'}
        >
          { accountName }
        </Text>
        <Text
          whiteSpace={"nowrap"}
          textOverflow={"ellipsis"}
          overflow={'hidden'}
        >
          { hoursUpdatedAgo } Hours Ago
        </Text>
      </Stack>


      <Box
        display={"flex"}
        alignItems={"flex-start"}
        justifyContent={"flex-start"}
      >
        <Badge colorScheme={isProfit ? "green" : "red"} borderRadius={"full"}>
          {percentageChange}%
        </Badge>
      </Box>


      <Box
        display={"flex"}
        alignItems={"flex-end"}
        justifyContent={"flex-end"}
      >
        <Text>${amount}</Text>
      </Box>

    </HStack>
  );
};

const AccountsTab: React.FC = () => {

  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");
  const { height } = useWindowDimensions();
  const [accountId, setAccountId] = useState(AccountCardData[0].id);

  const { isOpen, onOpen, onClose } = useDisclosure();
  const shouldOpenDrawer = useBreakpointValue({ base: true, lg: false});

  useEffect(() => {
    if (!shouldOpenDrawer) onClose();
  }, [shouldOpenDrawer, onClose]);

  const showAccount = (accountId: number) => {
    setAccountId(accountId);
    if (shouldOpenDrawer) onOpen();
  };

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <Box display={{ "lg": "none" }}>
        <Heading p="5px" fontSize={"3xl"}>Accounts</Heading>
        <Divider borderColor={dividerColor} />
      </Box>
      <Box
        h="full"
        w="full"
        display={"flex"}
        flex={"row"}
      >

        <Stack
          w="full"
          h="full"
          overflow={"auto"}
          maxH={height}
          spacing={"20px"}
          p={{ base: "10px", sm: "20px"  }}
        >

          <Stack spacing={"10px"}>
            <Text fontWeight={"bold"}>Credit Cards</Text>
            {AccountCardData.map((card, idx) => {
              return (
                <AccountCard
                  key={idx}
                  accountName={card.accountName}
                  hoursUpdatedAgo={card.hoursUpdatedAgo}
                  amount={card.amount}
                  percentageChange={card.percentageChange}
                  isProfit={card.isProfit}
                  onClick={() => { showAccount(card.id) }}
                />
              );
            })}
          </Stack>

          <Stack spacing={"10px"}>
            <Text fontWeight={"bold"}>Depository</Text>
            {AccountCardData.map((card, idx) => {
              return (
                <AccountCard
                  key={idx}
                  accountName={card.accountName}
                  hoursUpdatedAgo={card.hoursUpdatedAgo}
                  amount={card.amount}
                  percentageChange={card.percentageChange}
                  isProfit={card.isProfit}
                  onClick={() => { showAccount(card.id) }}
                />
              );
            })}
          </Stack>

        </Stack>

        <Divider
          w="10px"
          borderColor={dividerColor}
          orientation="vertical"
          display={{ base: "none", lg: "flex" }}
        />

        <Drawer
          isOpen={isOpen}
          placement="right"
          onClose={onClose}
          size={"full"}
        >
          <DrawerOverlay />
          <DrawerContent>
            <DrawerCloseButton />
            <DrawerBody
              backgroundColor={useColorModeValue("white", "black")}
              pt="50px"
            >
              <AccountDrawer accountId={accountId} />
            </DrawerBody>
          </DrawerContent>
        </Drawer>

        <Box
          w="full"
          display={{ base: "none", lg: "flex" }}
          p="20px"
        >
          <AccountDrawer accountId={accountId} />
        </Box>
      </Box>
    </Box>
  );
};

export default AccountsTab;
