import { useMutation, gql } from '@apollo/client';
import { Spacer, Stack } from '@chakra-ui/react';
import Button from './button';
import { Button as ChakraButton } from '@chakra-ui/react';
import { useRouter } from 'next/router';

const CREATE_STRIPE_PORTAL = gql`
  mutation CreateStripePortalSession($teamSlug: String!) {
    createPortalSession(teamSlug: $teamSlug) {
      sessionUrl
    }
  }
`;

interface StripeCheckoutFormProps {
  teamSlug: string;
  subscriptionActive?: boolean;

  onUpdateClick: () => any;
};

const StripeCheckoutForm: React.FC<StripeCheckoutFormProps> = ({ onUpdateClick, teamSlug, subscriptionActive }) => {

  const router = useRouter();
  const [createPortalSession, { loading: portalLoading, error: portalError }] = useMutation(CREATE_STRIPE_PORTAL);

  const handlePortalSession = async () => {
    const response = await createPortalSession({ variables: { teamSlug } });
    const sessionUrl = response.data?.createPortalSession?.sessionUrl;
    if (sessionUrl) {
      router.push(sessionUrl);
    } else {
      console.log('[stripe error]', portalError?.message);
    }
  };

  return (
    <Stack w="full" direction={{ base: "column", md: "row" }} overflow={"auto"}>

      <Button onClick={onUpdateClick} flip={true}>
        Update Subscription
      </Button>
      <Button flip={true} hidden={!subscriptionActive}>
        Buy Minutes
      </Button>
      <Button isDisabled={portalLoading} onClick={handlePortalSession} flip={true}>
        Manage Billing
      </Button>

      <Spacer />

      <ChakraButton
        colorScheme='red'
        isDisabled={portalLoading}
        onClick={handlePortalSession}
        hidden={!subscriptionActive}
      >
        Cancel Subscription
      </ChakraButton>

    </Stack>
  );
};

export default StripeCheckoutForm;
