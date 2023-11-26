import { Button as ChakraButton, ButtonProps, useColorModeValue } from "@chakra-ui/react";

interface NavButtonProps extends ButtonProps {
  flip?: boolean;
  gradient?: boolean;
}

const Button: React.FC<NavButtonProps> = (props) => {

  const bgColor = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

  return (
    <ChakraButton
      px={2}
      py={1}
      onClick={props.onClick}
      rounded={'md'}
      variant={"outline"}
      bgColor={!props.flip ? textColor : bgColor}
      textColor={!props.flip ? bgColor : textColor}
      size={"md"}
      _hover={{
        backgroundColor: !props.flip ? bgColor : textColor,
        textColor: !props.flip ? textColor : bgColor,
        bgGradient: props.gradient ? 'linear(to-tl, #007CF0, #01DFD8)' : '',
      }}
      { ...props }
    >
      {props.children}
    </ChakraButton>
  );

};

export default Button;
