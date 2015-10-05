<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 10.9.15
 * Time: 22:06
 */

namespace newsletters\Validators;

class CampaignValidator
{
    public function validateArrayOfEmails($attribute, $values, $parameters)
    {
        $checkEmails = function ($carry, $email) {
            return $carry && filter_var($email, FILTER_VALIDATE_EMAIL) !== false;
        };

        return (is_array($values) && array_reduce($values, $checkEmails, true));
    }
}
