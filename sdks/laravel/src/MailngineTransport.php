<?php

declare(strict_types=1);

namespace Mailngine;

use Mailngine\Exceptions\ApiException;
use Mailngine\Exceptions\MailngineException;
use Symfony\Component\Mailer\Envelope;
use Symfony\Component\Mailer\SentMessage;
use Symfony\Component\Mailer\Transport\TransportInterface;
use Symfony\Component\Mime\Address;
use Symfony\Component\Mime\Email;
use Symfony\Component\Mime\RawMessage;

/**
 * Symfony Mailer transport that sends emails through the Mailngine API.
 *
 * Integrates with Laravel's Mail system so you can use:
 *   Mail::mailer('mailngine')->to($user)->send(new WelcomeEmail());
 *
 * Configuration in config/mail.php:
 *   'mailers' => [
 *       'mailngine' => [
 *           'transport' => 'mailngine',
 *           'key' => env('MAILNGINE_API_KEY'),
 *       ],
 *   ],
 */
class MailngineTransport implements TransportInterface
{
    private Mailngine $client;

    public function __construct(Mailngine $client)
    {
        $this->client = $client;
    }

    /**
     * Send an email message via the Mailngine API.
     *
     * @throws MailngineException  If the API request fails.
     */
    public function send(RawMessage $message, ?Envelope $envelope = null): ?SentMessage
    {
        $envelope ??= Envelope::create($message);

        if (! $message instanceof Email) {
            throw new MailngineException(
                'MailngineTransport only supports Symfony\Component\Mime\Email instances'
            );
        }

        $from = $envelope->getSender();
        $to = array_map(fn (Address $a) => $a->getAddress(), $envelope->getRecipients());

        $params = [
            'from' => $from->getAddress(),
            'to' => $to,
            'subject' => $message->getSubject() ?? '',
        ];

        // CC
        $cc = $message->getCc();
        if (! empty($cc)) {
            $params['cc'] = array_map(fn (Address $a) => $a->getAddress(), $cc);
        }

        // BCC
        $bcc = $message->getBcc();
        if (! empty($bcc)) {
            $params['bcc'] = array_map(fn (Address $a) => $a->getAddress(), $bcc);
        }

        // Reply-To
        $replyTo = $message->getReplyTo();
        if (! empty($replyTo)) {
            $params['reply_to'] = $replyTo[0]->getAddress();
        }

        // HTML body
        $htmlBody = $message->getHtmlBody();
        if ($htmlBody !== null) {
            $params['html'] = (string) $htmlBody;
        }

        // Text body
        $textBody = $message->getTextBody();
        if ($textBody !== null) {
            $params['text'] = (string) $textBody;
        }

        // Custom headers
        $customHeaders = [];
        foreach ($message->getHeaders()->all() as $header) {
            $name = strtolower($header->getName());
            // Skip standard headers that are already handled as dedicated fields
            if (in_array($name, ['from', 'to', 'cc', 'bcc', 'reply-to', 'subject', 'content-type', 'mime-version', 'date', 'message-id'], true)) {
                continue;
            }
            $customHeaders[$header->getName()] = $header->getBodyAsString();
        }
        if (! empty($customHeaders)) {
            $params['headers'] = $customHeaders;
        }

        $this->client->emails()->send($params);

        return new SentMessage($message, $envelope);
    }

    public function __toString(): string
    {
        return 'mailngine';
    }
}
